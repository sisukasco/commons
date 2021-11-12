package redis_test

import (
	"strconv"
	"testing"
	"time"

	sredis "github.com/sisukasco/commons/redis"

	redigo "github.com/gomodule/redigo/redis"
	"syreclabs.com/go/faker"
)

const (
	redisURL = "redis://localhost:6379/5"
)

func TestSetGet(t *testing.T) {

	r, _ := sredis.New(redisURL)
	redisConn := r.Open()
	defer redisConn.Close()
	k := faker.RandomString(12)
	v := faker.RandomString(12)
	_, err := redisConn.Do("SET", k, v)
	if err != nil {
		t.Errorf("Error setting Redis value %v", err)
		return
	}
	res, err := redigo.String(redisConn.Do("GET", k))
	if err != nil {
		t.Errorf("Error getting Redis value %v", err)
		return
	}
	if v != res {
		t.Errorf("Redis Set does not match. expected %v vs %v", v, res)
	}
	t.Logf("Received value back from redis %s ", res)
}

func TestMultiDatabases(t *testing.T) {
	r1, _ := sredis.New("redis://localhost:6379/1")
	redisConn1 := r1.Open()
	defer redisConn1.Close()

	r2, _ := sredis.New("redis://localhost:6379/2")
	redisConn2 := r2.Open()
	defer redisConn2.Close()

	k := faker.RandomString(12)
	v1 := faker.RandomString(12)
	_, err := redisConn1.Do("SET", k, v1)
	if err != nil {
		t.Errorf("Error setting Redis value %v", err)
		return
	}

	v2 := faker.RandomString(12)
	_, err = redisConn2.Do("SET", k, v2)
	if err != nil {
		t.Errorf("Error setting Redis value %v", err)
		return
	}

	res1, err := redigo.String(redisConn1.Do("GET", k))
	if err != nil {
		t.Errorf("Error getting Redis value %v", err)
		return
	}
	if res1 != v1 {
		t.Errorf("Value didn't match. expected %v recvd %v", v1, res1)
	}
	t.Logf("k %s v1 %s res1 %s", k, v1, res1)

	res2, err := redigo.String(redisConn2.Do("GET", k))
	if err != nil {
		t.Errorf("Error getting Redis value %v", err)
		return
	}
	if res2 != v2 {
		t.Errorf("Value didn't match. expected %v recvd %v", v2, res2)
	}

	t.Logf("k %s v2 %s res2 %s", k, v2, res2)
}

//TODO: Get the redis URL from .env
func TestSortedSet(t *testing.T) {
	r, _ := sredis.New(redisURL)
	redisConn := r.Open()
	defer redisConn.Close()

	k := "recent_updated_forms"

	f1 := faker.RandomString(12)
	f2 := faker.RandomString(12)

	for i := 0; i < 5; i++ {
		redisConn.Do("ZINCRBY", k, 1, f1)
	}

	redisConn.Do("ZINCRBY", k, 1, f2)

	res, err := redigo.Strings(redisConn.Do("ZPOPMAX", k))

	if err != nil {
		t.Errorf("Error getting ZPOPMAX %v ", err)
		return
	}

	if res[0] != f1 {
		t.Errorf("ZPOPMAX expected %v got %v", f1, res[0])
	}

	score, err := strconv.Atoi(res[1])
	if err != nil {
		t.Errorf("Error strconv %v ", err)
		return
	}
	if score != 5 {
		t.Errorf("ZPOPMAX score does not match expedted 5, got %d ", score)
		return
	}
	t.Logf("ZPOPMAX result1 %v ", res)

	res2, err := redigo.Strings(redisConn.Do("ZPOPMAX", k))

	if err != nil {
		t.Errorf("Error getting ZPOPMAX %v ", err)
		return
	}
	t.Logf("ZPOPMAX result 2 %v ", res2)

	res3, err := redigo.Strings(redisConn.Do("ZPOPMAX", k))

	if err != nil {
		t.Errorf("Error getting ZPOPMAX %v ", err)
		return
	}
	t.Logf("ZPOPMAX result 3 %v ", res3)

	redisConn.Do("DEL", k)

}

func TestRoundRobinList(t *testing.T) {
	r, _ := sredis.New(redisURL)

	list_name := "test_couch_nodes"

	r.Delete(list_name)

	list := []int{1, 2, 3, 4, 5}
	N := len(list)

	err := r.CreateList(list_name, list)
	if err != nil {
		t.Fatalf("Error creating a list in redis %v", err)
	}
	for i := 0; i < 20; i++ {
		res, err := r.RoundRobinFromList(list_name)
		if err != nil {
			t.Fatalf("Error getting a list in redis %v", err)
		}
		t.Logf("Received value from list %#v", res)

		expected := N - (i - (i/N)*N)
		if res != expected {
			t.Errorf("TestRoundRobinList: Expected %v received %v", expected, res)
		}
	}

	//when 4
	//5 - 4

	// when 5
	//5 - 0(5-1*5)
	//when 6
	//5 - 1(6-(i/5)*5)
	//when 7
	//5 - 2(7-5)
	//when 8
	///5 - 3(8-5)
	// when 9
	///5 - 4(9-5)
	//when 10
	///5 - 0(10-2*5)
	//when 11
	//5 - 0(11-(i/5)*5)

	//5 - (i - (i/5)*5)
	//N - (i - (i/N)*N)

	r.Delete(list_name)
}

func TestHMap(t *testing.T) {

	r, _ := sredis.New(redisURL)

	key := "testmap" + faker.RandomString(12)

	t.Logf("TestHMap key %s", key)

	err := r.CreateHMap(key, map[string]string{"param1": "123"})
	if err != nil {
		t.Fatalf("Error creating a HMap in redis %v", err)
	}
	if !r.Exists(key) {
		t.Fatalf("Redis didn't create the key %s !", key)
	}

	v, err := r.GetHMapValue(key, "param1")
	if err != nil {
		t.Fatalf("Error getting a HMap in redis %v", err)
	}

	t.Logf("Received value %s", v)

	r.Delete(key)

	if r.Exists(key) {
		t.Fatalf("Redis didn't delete the key %s !", key)
	}
}
func TestSettingHMap(t *testing.T) {
	r, _ := sredis.New(redisURL)

	key := "testmap" + faker.RandomString(12)

	t.Logf("TestHMap key %s", key)

	err := r.CreateHMap(key, map[string]string{"param1": "123"})
	if err != nil {
		t.Fatalf("Error creating a HMap in redis %v", err)
	}

	newVal := "some456"
	err = r.SetHMapValue(key, "param2", newVal)
	if err != nil {
		t.Fatalf("Error setting a HMap in redis %v", err)
	}

	v, err := r.GetHMapValue(key, "param2")
	if err != nil {
		t.Fatalf("Error getting a HMap in redis %v", err)
	}

	if v != newVal {
		t.Fatalf("HMap setting value (param2) failed. Expected %s", newVal)
	}

}

func TestSetGetKeyValue(t *testing.T) {
	r, _ := sredis.New(redisURL)
	key := "test-key-set-" + faker.RandomString(12)
	value := "test-key-val-" + faker.RandomString(8)
	t.Logf("Key %s Value %s ", key, value)

	err := r.SetKeyValue(key, value, 2)
	if err != nil {
		t.Fatalf("Error setting key value in redis %v ", err)
	}
	v, err := r.GetKeyValue(key)
	if err != nil {
		t.Fatalf("Error getting key value from redis %v ", err)
	}
	if value != v {
		t.Errorf("Expected the same value %s %s", value, v)
	}
	t.Logf("Received Value %s ", value)

}

func TestGetKeyValueWithExpiry(t *testing.T) {
	r, _ := sredis.New(redisURL)
	key := "test-key-set-" + faker.RandomString(12)
	value := "test-key-val-" + faker.RandomString(8)
	t.Logf("Key %s Value %s ", key, value)

	err := r.SetKeyValue(key, value, 2)
	if err != nil {
		t.Fatalf("Error setting key value in redis %v ", err)
	}
	time.Sleep(3 * time.Second)
	_, err = r.GetKeyValue(key)
	if err == nil {
		t.Errorf("Expected that the key will expire in the specified time")
	}

}
