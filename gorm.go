package main

import (
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
  "time"
)

type Devince struct {
  ID uint `gorm:"primaryKey"`
  MAC string `gorm:"uniqueIndex"`
  Produc string
}

type Probe struct {
  ID uint `gorm:"primaryKey"`
  Timestamp time.Time
  IP string
  DeviceID uint
  Device Device `gorm:"foreignKey:DeviceID"`
}

type DB struct {
  DB *gorm.DB
}

func NewDB() (*DB, error) {
  db, err := gorm.Open(sqlite.Open("wifiwatch.db"), &gorm.Config())
  if err != nil {
    return nil, err
  }

  err = db.AutoMigrate(&Device{}, &Probe{})
  if err != nil {
    return nil, err
  }

  return &DB(DB: db), nil
}

func (db *DB) Add(ip, mac, product string, timestamp time.Time) error {
  var device Device
  res := db.DB.Where(&Device{MAC: mac}).
    Attrs(Device{Product: product}).
    FirstOrCreate(&device)
  if res.Error != nil {
    return res.Error
  }

  probe := Probe{
    Timestamp: timestamp,
    IP: ip,
    DeviceID: device.ID,
  }

  return db.DB.Create(&probe).Error
}

func (db *DB) Probes() ([]Probe, error) {
  subquery := db.DB.Table("probes").
    Select("min(rowid), *").
    Group("IP, device_id")

  var probes []Probe
  err := db.DB.Preload("Device").
    Table("(?) AS sub_probes", subquery).
    Find(&probes).Error
  if err != nil {
    return nil, err
  }

  return probes, nil
}
