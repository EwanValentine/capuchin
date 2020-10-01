use arrow::util::pretty;

use datafusion::datasource::csv::CsvReadOptions;
use datafusion::execution::context::ExecutionContext;

use std::thread;
use std::sync::{Arc, Mutex};
use std::error::Error;

use grpc::ServerRequestSingle;
use grpc::ServerHandlerContext;
use grpc::ServerResponseUnarySink;

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
    res: ServerResponseUnarySink<leader::leader::Response>) -> grpc::Result<()> {

    println!("Query: {:?}", req.message.query);

    let ctx = &mut self.ctx.lock().expect("could not lock");

    let plan = ctx.create_logical_plan(req.message.query.as_ref()).unwrap();
    let plan = ctx.optimize(&plan).unwrap();
    let plan = ctx.create_physical_plan(&plan, 1024 * 1024).unwrap();
    let results = ctx.collect(plan.as_ref()).unwrap();

    pretty::print_batches(&results).unwrap();

    res.finish(leader::leader::Response {
      status: "Ok".to_string(),
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

  println!("Testing 123");

  let port = 50051;
  let mut server_builder = grpc::ServerBuilder::new_plain();
  // server_builder.http.set_addr(&"127.0.0.1")?;
  server_builder.http.set_port(port);
  server_builder.add_service(leader::leader_grpc::LeaderAPIServer::new_service_def(svc));
 
  println!("Testing 456");

  let server = server_builder.build().expect("server");

  println!("server stared on addr {}", server.local_addr());

  loop {
    thread::park();
  }
}
