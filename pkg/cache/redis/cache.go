package redis

import (
    "log"
    "strconv"

    "github.com/gomodule/redigo/redis"
)

type Cache struct {
    conn redis.Conn
}

/**
 * Create a new cache by connecting to the local Redis server
 */
func New() (*Cache, error) {
    log.Println("creating new Redis cache...")
    conn, err := redis.DialURL("redis://localhost")
    if err != nil {
        return nil, err
    }
    return &Cache{conn}, nil
}

/**
 * Get int value from cache for given key
 */
func (c *Cache) GetInt(key interface{}) (int, error) {
    response, err := redis.Int(c.conn.Do("GET", key))
    if err != nil {
        return 0, err
    }
    return response, nil
}

/**
 * Get string value from cache for given key
 */
func (c *Cache) GetString(key interface{}) (string, error) {
    response, err := redis.String(c.conn.Do("GET", key))
    if err != nil {
        return "", err
    }
    return response, nil
}

/**
 * Set key-value pair in cache
 */
func (c *Cache) Set(key, value interface{}) error {
    _, err := c.conn.Do("SET", key, value)
    return err
}

/**
 * Set key-value pair in cache with expiration lifetime
 */
func (c *Cache) SetEx(key, value interface{}, lifetime int) error {
    // convert lifetime to string
    lifetimeString := strconv.Itoa(lifetime)
    log.Println("setting key-value pair with expiration in Redis cache...")
    _, err := c.conn.Do("SETEX", key, lifetimeString, value)
    if err != nil {
        log.Println("failed to set key-value pair in Redis cache")
        return err
    }
    log.Println("successfully stored key-value pair in Redis cache")
    return nil
}

/**
 * Delete key-value pair from cache
 */
func (c *Cache) Delete(key interface{}) error {
    _, err := c.conn.Do("DEL", key)
    return err
}
