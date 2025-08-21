package enum

import (
    "fmt"
    "slices"
    "strings"
)

// This file has been created automatically by `go-enum-generate`
// DO NOT MODIFY NOR EDIT THIS FILE DIRECTLY.
// To modify this enum, edit the enums.json or enums.yaml definition file
// To know more about `go-enum-generate`, see go to `https://github.com/debarbarinantoine/go-enum-generate`
// Generated at: 2025-08-21 20:58:55

type SessionStore uint

const (
    inMemory SessionStore = iota
    redis
    gORM
    mySQL
    postgreSQL
    mSSQL
    sQLite3
    mongoDB
)

var sessionStoreKeys = make(map[SessionStore]struct{}, 8)
var sessionStoreValues = make(map[string]SessionStore, 8)
var sessionStoreKeysArray = make([]SessionStore, 8)
var sessionStoreValuesArray = make([]string, 8)

func init() {
    sessionStoreKeys[inMemory] = struct{}{}
    sessionStoreKeysArray[0] = inMemory
    sessionStoreValues["memstore"] = inMemory
    sessionStoreValuesArray[0] = "memstore"

    sessionStoreKeys[redis] = struct{}{}
    sessionStoreKeysArray[1] = redis
    sessionStoreValues["redisstore"] = redis
    sessionStoreValuesArray[1] = "redisstore"

    sessionStoreKeys[gORM] = struct{}{}
    sessionStoreKeysArray[2] = gORM
    sessionStoreValues["gormstore"] = gORM
    sessionStoreValuesArray[2] = "gormstore"

    sessionStoreKeys[mySQL] = struct{}{}
    sessionStoreKeysArray[3] = mySQL
    sessionStoreValues["mysqlstore"] = mySQL
    sessionStoreValuesArray[3] = "mysqlstore"

    sessionStoreKeys[postgreSQL] = struct{}{}
    sessionStoreKeysArray[4] = postgreSQL
    sessionStoreValues["postgresstore"] = postgreSQL
    sessionStoreValuesArray[4] = "postgresstore"

    sessionStoreKeys[mSSQL] = struct{}{}
    sessionStoreKeysArray[5] = mSSQL
    sessionStoreValues["mssqlstore"] = mSSQL
    sessionStoreValuesArray[5] = "mssqlstore"

    sessionStoreKeys[sQLite3] = struct{}{}
    sessionStoreKeysArray[6] = sQLite3
    sessionStoreValues["sqlite3store"] = sQLite3
    sessionStoreValuesArray[6] = "sqlite3store"

    sessionStoreKeys[mongoDB] = struct{}{}
    sessionStoreKeysArray[7] = mongoDB
    sessionStoreValues["mongodbstore"] = mongoDB
    sessionStoreValuesArray[7] = "mongodbstore"
}

func (e SessionStore) String() string {
    switch e {
        case inMemory:
            return "memstore"
        case redis:
            return "redisstore"
        case gORM:
            return "gormstore"
        case mySQL:
            return "mysqlstore"
        case postgreSQL:
            return "postgresstore"
        case mSSQL:
            return "mssqlstore"
        case sQLite3:
            return "sqlite3store"
        case mongoDB:
            return "mongodbstore"
        default:
            return fmt.Sprintf("Unknown SessionStore (%d)", e.Value())
    }
}

func (e *SessionStore) Parse(str string) error {

    str = strings.TrimSpace(str)

    if val, ok := sessionStoreValues[str]; ok {
        *e = val
        return nil
    }
    return fmt.Errorf("invalid SessionStore: %s", str)
}

func (e SessionStore) Value() uint {
    return uint(e)
}

func (e SessionStore) MarshalText() ([]byte, error) {
    return []byte(e.String()), nil
}

func (e *SessionStore) UnmarshalText(text []byte) error {
    return e.Parse(string(text))
}

func (e SessionStore) IsValid() bool {
    if _, ok := sessionStoreKeys[e]; !ok {
        return false
    }
    return true
}

type sessionStores struct {
    InMemory SessionStore
    Redis SessionStore
    GORM SessionStore
    MySQL SessionStore
    PostgreSQL SessionStore
    MSSQL SessionStore
    SQLite3 SessionStore
    MongoDB SessionStore
}

var SessionStores = sessionStores{
    InMemory: inMemory,
    Redis: redis,
    GORM: gORM,
    MySQL: mySQL,
    PostgreSQL: postgreSQL,
    MSSQL: mSSQL,
    SQLite3: sQLite3,
    MongoDB: mongoDB,
}

func (e sessionStores) Values() []SessionStore {
    return slices.Clone(sessionStoreKeysArray)
}

func (e sessionStores) Args() []string {
    return slices.Clone(sessionStoreValuesArray)
}

func (e sessionStores) Description() string {
    var strBuilder strings.Builder
    strBuilder.WriteString("\tAvailable SessionStores:\n")
    for _, enumVal := range e.Values() {
        strBuilder.WriteString(fmt.Sprintf("=> %d -> %s\n", enumVal.Value(), enumVal.String()))
    }
    return strBuilder.String()
}

func (e sessionStores) Cast(value uint) (SessionStore, error) {
    if _, ok := sessionStoreKeys[SessionStore(value)]; !ok {
        return 0, fmt.Errorf("invalid cast SessionStore: %d", value)
    }
    return SessionStore(value), nil
}
