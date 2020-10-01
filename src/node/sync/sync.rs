extern crate csv;
extern crate rusoto_core;
extern crate rusoto_s3;

use futures::future;
use futures::TryStreamExt;

use std::fs::File;
use std::io::Write;
use bytes::buf::BufMutExt;
use std::default::Default;
use std::fmt::{self, Debug, Display, Formatter};

use rusoto_core::request::HttpClient;
use rusoto_core::credential::DefaultCredentialsProvider;
use rusoto_s3::*;
use rusoto_core::Region;
use rusoto_core::credential::ProfileProvider;

pub struct S3Sync {
  s3_client: rusoto_s3::S3Client,
}

pub struct SyncError(String);

impl Debug for SyncError {
  fn fmt(&self, f: &mut Formatter) -> fmt::Result {
    Display::fmt(self, f)
  }
}

impl Display for SyncError {
  fn fmt(&self, f: &mut Formatter) -> fmt::Result {
    write!(f, "{}", self.0)
  }
}

pub type SyncResult<T> = Result<T, SyncError>;

impl S3Sync {
  pub fn new() -> S3Sync {
    let provider = ProfileProvider::with_configuration(
      "/Users/ewanvalentine/.aws/credentials".to_string(),
      "ewanvalentine".to_string(),
    );

    let s3_client = S3Client::new_with(
      HttpClient::new().expect("failed to create request dispatcher"), 
      provider,
      Region::EuWest1
    );

    S3Sync {
      s3_client,
    }
  }

  pub async fn sync(&self, key: &str, bucket: &str) -> SyncResult<()> {
    let get_object_request = GetObjectRequest{
      key: key.to_string(),
      bucket: bucket.to_string(),
     ..Default::default()
   };

   let res = self.s3_client.get_object(get_object_request).await;
   
   let mut f = File::create("data.csv").unwrap();
   let stream = res.unwrap().body.unwrap();

   let body = stream.map_ok(|b| bytes::BytesMut::from(&b[..])).try_concat().await.unwrap();

    f.write(&body[..]).unwrap();
    Ok(())
  }

  pub async fn sync_list(&self, prefix: &str, source: &str, start: &str, end: &str) -> SyncResult<()> {

    // Load range from S3 into memory/arrow
    let list_objects_input = ListObjectsV2Request{
      bucket: source.to_string(),
      prefix: Some(prefix.to_string()),
      ..Default::default()
    };

    let results = self.s3_client.list_objects_v2(list_objects_input)
      .await
      .unwrap()
      .contents
      .unwrap_or_default()
      .into_iter()
      .collect::<Vec<_>>();

    let objects = results.into_iter()
      .map(|res| {
        let get_object_request = GetObjectRequest{
           key: res.key.unwrap().to_string(),
          bucket: source.to_string(),
          ..Default::default()
        };

        self.s3_client.get_object(get_object_request)
      });

    let content = future::join_all(objects).await;

    let mut f = File::create("data.csv").unwrap();

    let mut output = vec![].writer();
    for file in content.into_iter() {
      let stream = file.unwrap().body.unwrap();
      
      // Map into bytes
      let body = stream.map_ok(|b| bytes::BytesMut::from(&b[..])).try_concat().await.unwrap();
      output.write(&body[..]);
    }

    output.flush().unwrap();
    f.write_all(&output.into_inner());

    Ok(())
  }
}

// pub fn sync_example<T>(syncer: &T) where T: Syncer {}

#[cfg(test)]
mod tests {
  use super::*;

  #[test]
  fn can_write_to_local() {
    let s3_sync = S3Sync::new();
    s3_sync.sync(
      &"/ewans-amazing-test/".to_string(),
      &"kilimanjaro-latest-datalake".to_string(),
    );
  }
}