package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	// connect to the cluster
	cluster := gocql.NewCluster("172.26.0.2", "172.26.0.3", "172.26.0.4") //replace PublicIP with the IP addresses used by your cluster.
	cluster.Consistency = gocql.One
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "Username", Password: "Password", AllowedAuthenticators: []string{"com.instaclustr.cassandra.auth.InstaclustrPasswordAuthenticator"}} //replace the username and password fields with their real settings, you will need to allow the use of the Instaclustr Password Authenticator.
	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return
	}
	defer session.Close()

	// create keyspaces
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS sleep_centre WITH REPLICATION = {'class' : 'NetworkTopologyStrategy', 'datacenter1' : 3};").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	// create table
	err = session.Query("CREATE TABLE IF NOT EXISTS sleep_centre.sleep_study (name text, study_date date, sleep_time_hours float, PRIMARY KEY (name, study_date));").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	// insert some practice data
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-07', 8.2);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-08', 6.4);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-09', 7.5);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-07', 6.6);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-08', 6.3);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-09', 6.7);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Emily', '2018-01-07', 7.2);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Emily', '2018-01-09', 7.5);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Tom', '2018-01-09', 7.5);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Tomy', '2018-01-09', 8.5);").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	// Return average sleep time for James
	var sleep_time_hours float32

	sleep_time_output := session.Query("SELECT avg(sleep_time_hours) FROM sleep_centre.sleep_study WHERE name = 'Tom';").Iter()
	sleep_time_output.Scan(&sleep_time_hours)
	fmt.Println("Average sleep time for James was: ", sleep_time_hours, "h")

	// return average sleep time for group
	sleep_time_output = session.Query("SELECT avg(sleep_time_hours) FROM sleep_centre.sleep_study;").Iter()
	sleep_time_output.Scan(&sleep_time_hours)
	fmt.Println("Average sleep time for the group was: ", sleep_time_hours, "h")
}
