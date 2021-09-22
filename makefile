benchmark:
		@go test -bench=./... | tee ./graphic/out.dat ; \
			awk '/Benchmark/{count ++; gsub(/BenchmarkTest/,""); printf("%d,%s,%s,%s\n",count,$$1,$$2,$$3)}' ./graphic/out.dat > ./graphic/final.dat ; \
			gnuplot -e "file_path='./graphic/final.dat'" -e "graphic_file_name='./graphic/operations.png'" -e "y_label='number of operations'" -e "y_range_min='000000000''" -e "y_range_max='400000000'" -e "column_1=1" -e "column_2=3" ./graphic/performance.gp ; \
			gnuplot -e "file_path='./graphic/final.dat'" -e "graphic_file_name='./graphic/time_operations.png'" -e "y_label='each operation in nanoseconds'" -e "y_range_min='000''" -e "y_range_max='45000'" -e "column_1=1" -e "column_2=4" ./graphic/performance.gp ; \
			rm -f ./graphic/out.dat ./graphic/final.dat ; \
    
		echo "'graphic/operations.png' and 'graphic/time_operations.png' graphics were generated."

dependencies: 
	docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
		--name etcd quay.io/coreos/etcd:v2.3.8 \
		-name etcd0 \
		-advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 \
		-listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
		-initial-advertise-peer-urls http://${HostIP}:2380 \
		-listen-peer-urls http://0.0.0.0:2380 \
		-initial-cluster-token etcd-cluster-1 \
		-initial-cluster etcd0=http://${HostIP}:2380 \
		-initial-cluster-state new
