docker compose -f docker-compose.etcd.yml up -d
go test ./internal/storage -run '^$' -bench '^BenchmarkVClock' -benchmem -benchtime=5s -count=3 | tee bench_vclock.txt
python3 -m pip install pandas matplotlib seaborn
python3 plots.py