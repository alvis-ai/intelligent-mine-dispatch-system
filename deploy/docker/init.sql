CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    password VARCHAR(256) NOT NULL,
    real_name VARCHAR(64) DEFAULT '',
    email VARCHAR(128) DEFAULT '',
    phone VARCHAR(20) DEFAULT '',
    role INT DEFAULT 0,
    status INT DEFAULT 1,
    mine_id BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vehicles (
    id BIGINT PRIMARY KEY,
    plate VARCHAR(32) UNIQUE NOT NULL,
    type INT DEFAULT 1,
    status INT DEFAULT 1,
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,
    fuel_level DOUBLE PRECISION DEFAULT 100,
    mine_id BIGINT DEFAULT 0,
    driver_id BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dispatch_tasks (
    id BIGINT PRIMARY KEY,
    vehicle_id BIGINT NOT NULL,
    load_point_id BIGINT NOT NULL,
    dump_point_id BIGINT NOT NULL,
    material VARCHAR(64) DEFAULT '',
    load_lat DOUBLE PRECISION DEFAULT 0,
    load_lon DOUBLE PRECISION DEFAULT 0,
    dump_lat DOUBLE PRECISION DEFAULT 0,
    dump_lon DOUBLE PRECISION DEFAULT 0,
    status VARCHAR(32) DEFAULT 'pending',
    algorithm VARCHAR(32) DEFAULT 'fifo',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_mine_id ON users(mine_id);
CREATE INDEX idx_vehicles_mine_id ON vehicles(mine_id);
CREATE INDEX idx_vehicles_status ON vehicles(status);
CREATE INDEX idx_dispatch_tasks_vehicle_id ON dispatch_tasks(vehicle_id);
CREATE INDEX idx_dispatch_tasks_status ON dispatch_tasks(status);

CREATE TABLE IF NOT EXISTS vehicle_types (
    id BIGINT PRIMARY KEY,
    name VARCHAR(64) UNIQUE NOT NULL,
    description VARCHAR(256) DEFAULT '',
    icon VARCHAR(64) DEFAULT '',
    capacity DOUBLE PRECISION DEFAULT 0,
    weight DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS loading_points (
    id BIGINT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(32) NOT NULL DEFAULT 'loading',
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,
    material VARCHAR(64) DEFAULT '',
    status INT DEFAULT 1,
    mine_id BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_vehicle_types_name ON vehicle_types(name);
CREATE INDEX idx_loading_points_mine_id ON loading_points(mine_id);
CREATE INDEX idx_loading_points_type ON loading_points(type);

INSERT INTO vehicle_types (id, name, description, capacity, weight) VALUES
    (1, '矿用卡车', '大型矿山运输卡车', 60, 20),
    (2, '挖掘机', '矿山挖掘设备', 0, 30),
    (3, '装载机', '物料装载设备', 5, 15),
    (4, '推土机', '矿山推土设备', 0, 25)
ON CONFLICT (id) DO NOTHING;

INSERT INTO loading_points (id, name, type, material) VALUES
    (1, '装载点A', 'loading', '矿石'),
    (2, '装载点B', 'loading', '岩石'),
    (3, '卸载点C', 'dumping', '废石'),
    (4, '卸载点D', 'dumping', '矿石'),
    (5, '卸载点E', 'dumping', '岩石')
ON CONFLICT (id) DO NOTHING;
