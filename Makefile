all:
	@go run main.go;\

createDb:
	@cd ./test_suite && python3 populate.py;\

publishResult:
	@cd ./test_suite && python3 publish.py;\

deleteDb:
	@psql puppylove -c "delete from users;";\
	psql puppylove -c "delete from send_hearts;";\
	psql puppylove -c "delete from match_tables;";\
	psql puppylove -c "delete from heart_claims;";\
	psql puppylove -c "delete from return_hearts;";\

showAuth:
	@psql puppylove -c "select id,auth_c from users;";\