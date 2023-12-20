run:
	@go run main.go;\

createdb:
	@cd Stress\ test/;\
	echo -e "adminID\nadminPASS" | python3 populate.py;\
	cd ..;\

deletedb:
	@psql puppylove -c "delete from users;";\
	psql puppylove -c "delete from send_hearts;";\
	psql puppylove -c "delete from match_tables;";\
	psql puppylove -c "delete from heart_claims;";\
	psql puppylove -c "delete from return_hearts;";\