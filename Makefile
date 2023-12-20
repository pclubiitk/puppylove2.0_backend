run:
	@go run main.go;\

delete:
	@psql puppylove -c "delete from users;";\
	psql puppylove -c "delete from send_hearts;";\
	psql puppylove -c "delete from match_tables;";\
	psql puppylove -c "delete from heart_claims;";\
	psql puppylove -c "delete from return_hearts;";\