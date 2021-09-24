benchmark:
		@go test -bench=./... | tee ./graphic/out.dat ; \
			awk '/Benchmark/{count ++; gsub(/BenchmarkTest/,""); printf("%d,%s,%s,%s\n",count,$$1,$$2,$$3)}' ./graphic/out.dat > ./graphic/final.dat ; \
			gnuplot -e "file_path='./graphic/final.dat'" -e "graphic_file_name='./graphic/operations.png'" -e "y_label='number of operations'" -e "y_range_min='000000000''" -e "y_range_max='400000000'" -e "column_1=1" -e "column_2=3" ./graphic/performance.gp ; \
			gnuplot -e "file_path='./graphic/final.dat'" -e "graphic_file_name='./graphic/time_operations.png'" -e "y_label='each operation in nanoseconds'" -e "y_range_min='000''" -e "y_range_max='45000'" -e "column_1=1" -e "column_2=4" ./graphic/performance.gp ; \
			rm -f ./graphic/out.dat ./graphic/final.dat ; \
    
		echo "'graphic/operations.png' and 'graphic/time_operations.png' graphics were generated."

dependencies:
		rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
				docker rmi gcr.io/etcd-development/etcd:v3.5.0 || true && \
				docker run \
				-p 2379:2379 \
				-p 2380:2380 \
				--mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
				--name etcd-gcr-v3.5.0 \
				gcr.io/etcd-development/etcd:v3.5.0 \
				/usr/local/bin/etcd \
				--name s1 \
				--data-dir /etcd-data \
				--listen-client-urls http://0.0.0.0:2379 \
				--advertise-client-urls http://0.0.0.0:2379 \
				--listen-peer-urls http://0.0.0.0:2380 \
				--initial-advertise-peer-urls http://0.0.0.0:2380 \
				--initial-cluster s1=http://0.0.0.0:2380 \
				--initial-cluster-token tkn \
				--initial-cluster-state new \
				--log-level info \
				--logger zap \
				--log-outputs stderr
