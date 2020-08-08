use std::env;
use std::thread;

use futures::executor;

mod leader;

use leader::leader::{Request, Response};
use leader::leader_grpc;
use crate::leader::leader_grpc::LeaderAPIClient;
use grpc::ClientStubExt;

fn main() {
    let port = 50051;
    let client_conf = Default::default();
    let client = leader_grpc::LeaderAPIClient::new_plain("::1", port, client_conf).unwrap();

    let mut req =  Request::new();
    req.set_query("SELECT * FROM test LIMIT 100".parse().unwrap());

    let resp = client
        .query(grpc::RequestOptions::new(), req)
        .join_metadata_result();

    println!("{:?}", executor::block_on(resp));
}