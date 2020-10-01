use serde_json;

use datafusion::datasource::csv::CsvReadOptions;
use datafusion::execution::context::ExecutionContext;
use datafusion::prelude::*;

use std::thread;
use std::sync::{Arc, Mutex};
use std::error::Error;

use grpc::ServerRequestSingle;
use grpc::ServerHandlerContext;
use grpc::ServerResponseUnarySink;
use grpc::ServerResponseSink;

mod leader;
mod sync;

use crate::leader::leader_grpc::LeaderAPI;
use crate::sync::sync::S3Sync;

use futures::executor::block_on;

type QueryContext = datafusion::execution::context::ExecutionContext;

struct LeaderAPIImpl {
  ctx: Arc<Mutex<QueryContext>>,
}

impl LeaderAPI {
  fn new_with_arrow(ctx: Arc<Mutex<QueryContext>>) -> LeaderAPIImpl {
    LeaderAPIImpl{
      ctx,
    }
  }
}

impl LeaderAPI for LeaderAPIImpl {
  fn query(
    &self,
    _o: ServerHandlerContext,
    req: ServerRequestSingle<leader::leader::Request>,
    res: ServerResponseSink<leader::leader::Response>) -> grpc::Result<()> {

    println!("Query: {:?}", req.message.query);

    let ctx = &mut self.ctx.lock().expect("could not lock");

    // register csv file with the execution context
    ctx.register_csv(
        "test",
        &"data.csv",
        CsvReadOptions::new(),
    ).unwrap();

    // execute the query
    let df = ctx.sql(&req.message.query, 10000).unwrap();

    let results = df.collect().await.unwrap();
    let serialized = serde_json::to_vec(&results).unwrap();

    res.finish(leader::leader::Response {
      status: "Ok".to_string(),
      data: serialized.as_bytes(),
      ..Default::default()
    })
  }
}

#[tokio::main]
async fn main() -> std::result::Result<(), Box<dyn Error>> {
  let s3_syncer = S3Sync::new();
  s3_syncer.sync(&"data.csv", &"ewanstestdata").await.unwrap();

  let mut ctx = ExecutionContext::new();
  ctx.register_csv(
    "test",
    "data.csv",
    CsvReadOptions::new(),
  ).unwrap();

  let safe_ctx = Arc::new(Mutex::new(ctx));
  let svc = LeaderAPI::new_with_arrow(safe_ctx);

  let port = 50051;
  let mut server_builder = grpc::ServerBuilder::new_plain();
  // server_builder.http.set_addr(&"127.0.0.1")?;
  server_builder.http.set_port(port);
  server_builder.add_service(leader::leader_grpc::LeaderAPIServer::new_service_def(svc));
 
  let server = server_builder.build().expect("server");

  println!("server stared on addr {}", server.local_addr());

  loop {
    thread::park();
  }
}
