use arrow::util::pretty;

use datafusion::error::Result;
use datafusion::datasource::csv::CsvReadOptions;
use datafusion::execution::context::ExecutionContext;

use std::thread;
use std::sync::{Arc, Mutex};

use futures::stream::StreamExt;

use grpc::Metadata;
use grpc::ServerRequest;
use grpc::ServerResponseSink;
use grpc::ServerRequestSingle;
use grpc::ServerHandlerContext;
use grpc::ServerResponseUnarySink;

mod leader;

use leader::leader_grpc;
use crate::leader::leader_grpc::LeaderAPI;
use std::ops::{Deref, DerefMut};

type QueryContext = datafusion::execution::context::ExecutionContext;

struct LeaderAPIImpl {
  ctx: Arc<Mutex<QueryContext>>
}

impl dyn LeaderAPI {
  pub fn new_with_arrow(ctx: Arc<Mutex<QueryContext>>) -> LeaderAPIImpl {
    LeaderAPIImpl{
      ctx,
    }
  }
}

impl LeaderAPI for LeaderAPIImpl{
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

fn main() {

  let mut ctx = ExecutionContext::new();
  ctx.register_csv(
    "test",
    "data/train.csv",
    CsvReadOptions::new(),
  ).unwrap();

  let safe_ctx = Arc::new(Mutex::new(ctx));
  let svc = LeaderAPI::new_with_arrow(safe_ctx);

  let port = 50051;
  let mut server_builder = grpc::ServerBuilder::new_plain();
  server_builder.http.set_addr(String::from("localhost"));
  server_builder.http.set_port(port);
  server_builder.add_service(leader::leader_grpc::LeaderAPIServer::new_service_def(svc));
  
  let server = server_builder.build().expect("server");

  println!("server stared on addr {}", server.local_addr());

  loop {
    thread::park();
  }
}
