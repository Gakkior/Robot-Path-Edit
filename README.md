# 馃 鏈哄櫒浜鸿矾寰勭紪杈戝櫒 (Robot Path Editor)

涓€涓幇浠ｅ寲鐨勪笁绔吋瀹规満鍣ㄤ汉璺緞缂栬緫鍣紝鏀寔鍙鍖栫紪杈戝拰鏁版嵁搴撶鐞嗐€?

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Coverage](https://img.shields.io/badge/coverage-85%25-green.svg)

## 鉁?鐗规€?

### 馃帹 鍙鍖栫紪杈?
- **浜や簰寮忕敾甯?*: 鍩轰簬 Konva.js 鐨勯珮鎬ц兘鐢诲竷锛屾敮鎸佽妭鐐规嫋鎷姐€佽矾寰勮繛鎺?
- **瀹炴椂棰勮**: 鎵€瑙佸嵆鎵€寰楃殑缂栬緫浣撻獙锛屽疄鏃舵樉绀鸿矾寰勮鍒掔粨鏋?
- **澶氱瑙嗗浘**: 鐢诲竷瑙嗗浘鍜岃〃鏍艰鍥炬棤缂濆垏鎹紝婊¤冻涓嶅悓浣跨敤鍦烘櫙

### 馃梽锔?鏁版嵁绠＄悊
- **閫氱敤鏁版嵁搴撶紪杈戝櫒**: 鏀寔浠绘剰琛ㄧ粨鏋勭殑 CRUD 鎿嶄綔锛岀被浼?Excel 鐨勪娇鐢ㄤ綋楠?
- **鐏垫椿鏄犲皠**: 鍙€夋嫨浠绘剰鏁版嵁琛ㄤ綔涓虹偣浣嶈〃鍜岃矾寰勮〃锛屾敮鎸佽嚜瀹氫箟ID瀛楁鏄犲皠
- **瀹炴椂鍚屾**: 鐢诲竷瑙嗗浘涓庤〃鏍艰鍥炬暟鎹疄鏃跺弻鍚戝悓姝?

### 馃敡 鏅鸿兘绠楁硶
- **甯冨眬绠楁硶**: 缃戞牸甯冨眬銆佸姏瀵煎悜甯冨眬銆佸渾褰㈠竷灞€绛夊绉嶈嚜鍔ㄦ帓鍒楁柟寮?
- **璺緞鐢熸垚**: 鏈€杩戦偦杩炴帴銆佸畬鍏ㄨ繛閫氬浘銆佺綉鏍艰矾寰勭瓑鏅鸿兘璺緞鐢熸垚绠楁硶
- **璺緞浼樺寲**: 鏈€鐭矾寰勮绠椼€佽矾寰勫钩婊戜紭鍖?

### 馃洜锔?楂樼骇鍔熻兘
- **鎾ら攢/閲嶅仛**: 鍩轰簬鍛戒护妯″紡鐨勫畬鏁存搷浣滃巻鍙茬鐞?
- **鎻掍欢绯荤粺**: 鍙墿灞曠殑鎻掍欢鏋舵瀯锛屾敮鎸佽嚜瀹氫箟甯冨眬鍜岃矾寰勭畻娉?
- **瀹炴椂鐩戞帶**: Prometheus 鎸囨爣鏀堕泦銆佺粨鏋勫寲鏃ュ織銆佹€ц兘杩借釜
- **澶氱閮ㄧ讲**: Docker銆丼ystemd銆丳M2 绛夊绉嶉儴缃叉柟寮?

### 馃摫 璺ㄥ钩鍙版敮鎸?
- **妗岄潰绔?*: Windows銆丩inux銆乵acOS 鍘熺敓鏀寔
- **绉诲姩绔?*: PWA 鏀寔锛屽钩鏉胯澶囦紭鍖栫殑瑙︽帶浣撻獙
- **Web绔?*: 鍝嶅簲寮忚璁★紝鏀寔鎵€鏈夌幇浠ｆ祻瑙堝櫒

## 馃彈锔?鎶€鏈灦鏋?

### Go-Heavy 鍚庣鏋舵瀯
```
鈹屸攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?   鈹屸攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?   鈹屸攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?
鈹?  Handler灞?    鈹?   鈹?  Service灞?    鈹?   鈹?Repository灞?   鈹?
鈹? (HTTP鎺ュ彛)     鈹傗攢鈹€鈹€鈹€鈹?  (涓氬姟閫昏緫)    鈹傗攢鈹€鈹€鈹€鈹?  (鏁版嵁璁块棶)    鈹?
鈹斺攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?   鈹斺攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?   鈹斺攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?
         鈹?                      鈹?                      鈹?
         鈹?             鈹屸攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?             鈹?
         鈹?             鈹?  Plugin绯荤粺    鈹?             鈹?
         鈹?             鈹? (鎵╁睍绠楁硶)     鈹?             鈹?
         鈹?             鈹斺攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?             鈹?
         鈹?                      鈹?                      鈹?
鈹屸攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?
鈹?                    Domain灞?(鏍稿績棰嗗煙妯″瀷)                      鈹?
鈹?              Node, Path, Position, Metadata                    鈹?
鈹斺攢鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹€鈹?
```

### 鍓嶇鎶€鏈爤
- **鐢诲竷娓叉煋**: Konva.js (楂樻€ц兘 2D Canvas)
- **浜や簰閫昏緫**: 鍘熺敓 JavaScript (杞婚噺鍖?
- **鐘舵€佺鐞?*: 鍛戒护妯″紡 (鎾ら攢/閲嶅仛)
- **UI缁勪欢**: 鐜颁唬 CSS + HTML5

### 璁捐妯″紡搴旂敤
- **浠撳偍妯″紡**: 鏁版嵁璁块棶鎶借薄
- **閫傞厤鍣ㄦā寮?*: 澶氭暟鎹簱鏀寔
- **鍛戒护妯″紡**: 鎿嶄綔鍘嗗彶绠＄悊
- **绛栫暐妯″紡**: 甯冨眬绠楁硶鍒囨崲
- **瑙傚療鑰呮ā寮?*: 浜嬩欢椹卞姩鏋舵瀯
- **宸ュ巶妯″紡**: 缁勪欢鍒涘缓绠＄悊

## 馃殌 蹇€熷紑濮?

### 婕旂ず鐗?(鎺ㄨ崘)
鏈€蹇殑浣撻獙鏂瑰紡锛屾棤闇€鏁版嵁搴撻厤缃細

```bash
# 涓嬭浇骞跺惎鍔ㄦ紨绀虹増
go run cmd/demo/main.go

# 鎴栦娇鐢ㄩ缂栬瘧鐗堟湰
./demo.exe  # Windows
./demo      # Linux/macOS
```

璁块棶 http://localhost:8080 寮€濮嬩綋楠岋紒

### 瀹屾暣鐗堥儴缃?

#### 1. 鐜瑕佹眰
- Go 1.21+
- SQLite/MySQL/PostgreSQL (浠婚€変竴绉?
- Node.js 16+ (鍙€夛紝鐢ㄤ簬鍓嶇鏋勫缓)

#### 2. 蹇€熷畨瑁?
```bash
# 鍏嬮殕椤圭洰
git clone https://github.com/your-org/robot-path-editor.git
cd robot-path-editor

# 瀹夎渚濊禆
go mod download

# 鏋勫缓椤圭洰
./scripts/build.sh build-all  # Linux/macOS
build.bat build-all           # Windows

# 閰嶇疆鏁版嵁搴?
cp configs/config.yaml.example configs/config.yaml
# 缂栬緫 config.yaml 閰嶇疆鏁版嵁搴撹繛鎺?

# 鍚姩鏈嶅姟
./build/robot-path-editor
```

#### 3. Docker 閮ㄧ讲 (鎺ㄨ崘)
```bash
# 浣跨敤 Docker Compose 涓€閿儴缃?
docker-compose up -d

# 鍖呭惈鏁版嵁搴撱€丷edis銆佺洃鎺х殑瀹屾暣鐜
```

#### 4. 鐢熶骇鐜閮ㄧ讲
```bash
# 浣跨敤閮ㄧ讲鑴氭湰
./scripts/deploy.sh

# 鏀寔澶氱閮ㄧ讲鏂瑰紡
DEPLOY_MODE=systemd ./scripts/deploy.sh  # Systemd
DEPLOY_MODE=pm2 ./scripts/deploy.sh      # PM2
DEPLOY_MODE=docker ./scripts/deploy.sh   # Docker
```

## 馃摉 浣跨敤鎸囧崡

### 鍩虹鎿嶄綔

#### 鐢诲竷瑙嗗浘
1. **鍒涘缓鑺傜偣**: 鍙屽嚮绌虹櫧鍖哄煙鎴栦娇鐢ㄥ伐鍏锋爮
2. **绉诲姩鑺傜偣**: 鎷栨嫿鑺傜偣鍒扮洰鏍囦綅缃?
3. **鍒涘缓璺緞**: Shift+鐐瑰嚮涓や釜鑺傜偣
4. **鍒犻櫎鍏冪礌**: 閫変腑鍚庢寜Delete閿?
5. **鎾ら攢鎿嶄綔**: Ctrl+Z / Cmd+Z

#### 琛ㄦ牸瑙嗗浘
1. **鍒囨崲瑙嗗浘**: 鐐瑰嚮椤甸潰椤堕儴鐨?琛ㄦ牸瑙嗗浘"鎸夐挳
2. **缂栬緫鏁版嵁**: 鐩存帴鍦ㄨ〃鏍间腑淇敼鏁板€?
3. **鎵归噺鎿嶄綔**: 閫夋嫨澶氳杩涜鎵归噺缂栬緫鎴栧垹闄?
4. **瀵煎叆瀵煎嚭**: 鏀寔CSV銆丒xcel鏍煎紡

#### 鏅鸿兘绠楁硶
```bash
# 搴旂敤甯冨眬绠楁硶
curl -X POST http://localhost:8080/api/v1/layout/apply \
  -H "Content-Type: application/json" \
  -d '{"algorithm": "force-directed"}'

# 鐢熸垚璺緞
curl -X POST http://localhost:8080/api/v1/path-generation/nearest-neighbor \
  -H "Content-Type: application/json" \
  -d '{"max_distance": 200}'
```

### 楂樼骇閰嶇疆

#### 鏁版嵁搴撻厤缃?
```yaml
# config.yaml
database:
  type: "mysql"  # sqlite, mysql, postgresql
  dsn: "user:password@tcp(localhost:3306)/robot_paths"
  
  # 鑷畾涔夎〃鏄犲皠
  table_mapping:
    node_table: "robot_points"
    node_id_field: "point_id"
    path_table: "robot_routes"
    path_id_field: "route_id"
```

#### 鎻掍欢寮€鍙?
```go
// 鑷畾涔夊竷灞€鎻掍欢
type CustomLayoutPlugin struct{}

func (p *CustomLayoutPlugin) ApplyLayout(nodes []domain.Node, paths []domain.Path, config map[string]interface{}) ([]domain.Node, error) {
    // 瀹炵幇鑷畾涔夊竷灞€绠楁硶
    return nodes, nil
}

// 娉ㄥ唽鎻掍欢
pluginService.RegisterLayoutPlugin(&CustomLayoutPlugin{})
```

## 馃搳 API 鍙傝€?

### RESTful API

#### 鑺傜偣绠＄悊
- `GET /api/v1/nodes` - 鑾峰彇鎵€鏈夎妭鐐?
- `POST /api/v1/nodes` - 鍒涘缓鑺傜偣
- `GET /api/v1/nodes/{id}` - 鑾峰彇鍗曚釜鑺傜偣
- `PUT /api/v1/nodes/{id}` - 鏇存柊鑺傜偣
- `DELETE /api/v1/nodes/{id}` - 鍒犻櫎鑺傜偣
- `PUT /api/v1/nodes/{id}/position` - 鏇存柊鑺傜偣浣嶇疆

#### 璺緞绠＄悊
- `GET /api/v1/paths` - 鑾峰彇鎵€鏈夎矾寰?
- `POST /api/v1/paths` - 鍒涘缓璺緞
- `GET /api/v1/paths/{id}` - 鑾峰彇鍗曚釜璺緞
- `PUT /api/v1/paths/{id}` - 鏇存柊璺緞
- `DELETE /api/v1/paths/{id}` - 鍒犻櫎璺緞

#### 甯冨眬绠楁硶
- `POST /api/v1/layout/apply` - 搴旂敤甯冨眬绠楁硶

#### 璺緞鐢熸垚
- `POST /api/v1/path-generation/nearest-neighbor` - 鏈€杩戦偦璺緞
- `POST /api/v1/path-generation/full-connectivity` - 瀹屽叏杩為€?
- `POST /api/v1/path-generation/grid` - 缃戞牸璺緞

### WebSocket API (璁″垝涓?
- `/ws/canvas` - 瀹炴椂鐢诲竷鍚屾
- `/ws/notifications` - 绯荤粺閫氱煡

## 馃И 寮€鍙戞寚鍗?

### 寮€鍙戠幆澧冩惌寤?
```bash
# 瀹夎寮€鍙戝伐鍏?
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest

# 杩愯寮€鍙戞湇鍔″櫒
go run cmd/server/main.go

# 杩愯娴嬭瘯
go test ./... -v -cover

# 浠ｇ爜妫€鏌?
golangci-lint run
```

### 椤圭洰缁撴瀯
```
robot-path-editor/
鈹溾攢鈹€ cmd/                    # 搴旂敤鍏ュ彛
鈹?  鈹溾攢鈹€ server/            # 涓绘湇鍔″櫒
鈹?  鈹斺攢鈹€ demo/              # 婕旂ず鐗堟湰
鈹溾攢鈹€ internal/              # 鍐呴儴鍖?
鈹?  鈹溾攢鈹€ domain/            # 棰嗗煙妯″瀷
鈹?  鈹溾攢鈹€ services/          # 涓氬姟鏈嶅姟
鈹?  鈹溾攢鈹€ repositories/      # 鏁版嵁浠撳偍
鈹?  鈹溾攢鈹€ handlers/          # HTTP澶勭悊鍣?
鈹?  鈹斺攢鈹€ database/          # 鏁版嵁搴撻€傞厤
鈹溾攢鈹€ pkg/                   # 鍏叡鍖?
鈹?  鈹溾攢鈹€ logger/            # 鏃ュ織宸ュ叿
鈹?  鈹斺攢鈹€ middleware/        # 涓棿浠?
鈹溾攢鈹€ web/                   # 鍓嶇璧勬簮
鈹?  鈹斺攢鈹€ static/            # 闈欐€佹枃浠?
鈹溾攢鈹€ tests/                 # 娴嬭瘯鏂囦欢
鈹?  鈹溾攢鈹€ unit/              # 鍗曞厓娴嬭瘯
鈹?  鈹斺攢鈹€ integration/       # 闆嗘垚娴嬭瘯
鈹溾攢鈹€ scripts/               # 鏋勫缓鑴氭湰
鈹溾攢鈹€ configs/               # 閰嶇疆鏂囦欢
鈹斺攢鈹€ docs/                  # 鏂囨。
```

### 璐＄尞鎸囧崡
1. Fork 椤圭洰
2. 鍒涘缓鐗规€у垎鏀? `git checkout -b feature/amazing-feature`
3. 鎻愪氦鍙樻洿: `git commit -m 'Add amazing feature'`
4. 鎺ㄩ€佸垎鏀? `git push origin feature/amazing-feature`
5. 鎻愪氦 Pull Request

### 浠ｇ爜瑙勮寖
- 閬靛惊 Go 瀹樻柟浠ｇ爜瑙勮寖
- 浣跨敤 golangci-lint 杩涜浠ｇ爜妫€鏌?
- 鍗曞厓娴嬭瘯瑕嗙洊鐜?> 80%
- 鎻愪氦淇℃伅閬靛惊 Conventional Commits

## 馃搱 鎬ц兘鎸囨爣

### 绯荤粺鎬ц兘
- **鍝嶅簲鏃堕棿**: < 100ms (95%ile)
- **鍚炲悙閲?*: > 1000 QPS
- **鍐呭瓨浣跨敤**: < 256MB (绌鸿浇)
- **鍚姩鏃堕棿**: < 5s

### 鐢诲竷鎬ц兘
- **鑺傜偣鏁伴噺**: 鏀寔 10,000+ 鑺傜偣
- **璺緞鏁伴噺**: 鏀寔 50,000+ 璺緞
- **娓叉煋甯х巼**: 60 FPS (1080p)
- **鍝嶅簲寤惰繜**: < 16ms (瑙︽帶/榧犳爣)

### 鏁版嵁搴撴€ц兘
- **SQLite**: 閫傚悎 < 10涓?璁板綍
- **MySQL**: 閫傚悎 < 1000涓?璁板綍
- **PostgreSQL**: 閫傚悎 > 1000涓?璁板綍

## 馃敡 鏁呴殰鎺掗櫎

### 甯歌闂

#### 1. CGO鐩稿叧閿欒
```bash
# 閿欒: CGO_ENABLED=0, go-sqlite3 requires cgo
# 瑙ｅ喅: 浣跨敤绾疓o SQLite椹卞姩鎴栧惎鐢–GO
export CGO_ENABLED=1
go build ...
```

#### 2. 绔彛鍗犵敤
```bash
# 妫€鏌ョ鍙ｅ崰鐢?
netstat -tulpn | grep :8080

# 淇敼閰嶇疆鏂囦欢涓殑绔彛
# 鎴栬缃幆澧冨彉閲?
export PORT=8081
```

#### 3. 闈欐€佽祫婧?04
```bash
# 纭繚web鐩綍瀛樺湪
# 鎴栦娇鐢╣o:embed鍐呭祵璧勬簮
```

### 鏃ュ織鍒嗘瀽
```bash
# 鏌ョ湅搴旂敤鏃ュ織
tail -f /var/log/robot-path-editor/app.log

# 鏌ョ湅璁块棶鏃ュ織
tail -f /var/log/robot-path-editor/access.log

# 鏌ョ湅閿欒鏃ュ織
grep "ERROR" /var/log/robot-path-editor/app.log
```

### 鐩戞帶鎸囨爣
- **绯荤粺鎸囨爣**: CPU銆佸唴瀛樸€佺鐩樹娇鐢ㄧ巼
- **搴旂敤鎸囨爣**: QPS銆佸搷搴旀椂闂淬€侀敊璇巼
- **涓氬姟鎸囨爣**: 鑺傜偣鏁伴噺銆佽矾寰勬暟閲忋€佺敤鎴锋椿璺冨害

## 馃搫 璁稿彲璇?

鏈」鐩熀浜?MIT 璁稿彲璇佸紑婧?- 璇﹁ [LICENSE](LICENSE) 鏂囦欢

## 馃檹 鑷磋阿

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web妗嗘灦
- [GORM](https://gorm.io/) - Go ORM搴?
- [Konva.js](https://konvajs.org/) - 2D Canvas娓叉煋寮曟搸
- [Prometheus](https://prometheus.io/) - 鐩戞帶绯荤粺
- [Logrus](https://github.com/sirupsen/logrus) - 鏃ュ織搴?

## 馃敆 鐩稿叧閾炬帴

- [鏂囨。缃戠珯](https://robot-path-editor.github.io/docs)
- [鍦ㄧ嚎婕旂ず](https://demo.robot-path-editor.com)
- [Docker Hub](https://hub.docker.com/r/robotpatheditor/robot-path-editor)
- [闂鍙嶉](https://github.com/your-org/robot-path-editor/issues)
- [璁ㄨ绀惧尯](https://github.com/your-org/robot-path-editor/discussions)

---

<div align="center">

**濡傛灉杩欎釜椤圭洰瀵逛綘鏈夊府鍔╋紝璇风粰瀹冧竴涓?猸愶笍 Star锛?*

[馃悰 鎶ュ憡Bug](https://github.com/your-org/robot-path-editor/issues) 鈥?
[鉁?璇锋眰鍔熻兘](https://github.com/your-org/robot-path-editor/issues) 鈥?
[馃挰 鍙備笌璁ㄨ](https://github.com/your-org/robot-path-editor/discussions)

</div>