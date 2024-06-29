/*
 * This file was last modified at 2024-06-24 22:51 by Victor N. Skurikhin.
 * config_test.go
 * $Id$
 */

package env

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	testDataBaseDSN    = "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	testConfigFileName string
	testCryptoKey      []string
	testKey            string
	testServerAddress  string
	testTempFileName   string
)

func TestConfig(t *testing.T) {

	oldArgs := os.Args
	os.Args = []string{"agent", "-r"}
	getEnvironments()
	initServerFlags()
	reportInterval := 1 * time.Nanosecond
	flag.reportInterval = &reportInterval
	pollInterval := 2 * time.Nanosecond
	flag.pollInterval = &pollInterval
	getTestConfig()

	var tests = []struct {
		name string
		fCfg func() Config
		want string
	}{
		{
			name: "Test config #1",
			fCfg: GetAgentConfig,
			want: `
	dataBaseDSN     : 
	fileStoragePath : 
	key             : 
	pollInterval    : 2ns
	reportInterval  : 1ns
	restore         : false
	serverAddress   : localhost:8080
	storeInterval   : 0s
	urlHost         : http://localhost:8080
`,
		},
		{
			name: "Test config #2",
			fCfg: GetServerConfig,
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : /tmp/metrics-db.json
	key             : 
	pollInterval    : 0s
	reportInterval  : 0s
	restore         : true
	serverAddress   : localhost:8080
	storeInterval   : 5m0s
	urlHost         : http://localhost:8080
`,
		},
		{
			name: "Test config #3",
			fCfg: getTestConfig,
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + testTempFileName + `
	key             : ` + testKey + `
	pollInterval    : 1m0s
	reportInterval  : 1h0m0s
	restore         : true
	serverAddress   : ` + testServerAddress + `
	storeInterval   : 24h0m0s
	urlHost         : http://` + testServerAddress + `
`,
		},
		{
			name: "Test config #4",
			fCfg: func() Config { return getTestEnvironments(t) },
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + testTempFileName + `
	key             : ` + testKey + `
	pollInterval    : 0s
	reportInterval  : 0s
	restore         : true
	serverAddress   : ` + testServerAddress + `
	storeInterval   : 1s
	urlHost         : http://` + testServerAddress + `
`,
		},
		{
			name: "Test config #5",
			fCfg: getTestFlags,
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + testTempFileName + `
	key             : ` + testKey + `
	pollInterval    : 1m0s
	reportInterval  : 1h0m0s
	restore         : true
	serverAddress   : ` + testServerAddress + `
	storeInterval   : 24h0m0s
	urlHost         : http://` + testServerAddress + `
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			onceCfg = new(sync.Once)
			c := test.fCfg()
			got := c.String()
			assert.Equal(t, test.want, got)
		})
	}
	assert.True(t, cfg.IsDBSetup())
	os.Args = oldArgs
}

func TestTestFlags(t *testing.T) {
	cfg := getTestFlags()
	assert.Equal(t, testConfigFileName, cfg.ConfigFileName())
	assert.Equal(t, testCryptoKey, cfg.CryptoKey())
	assert.Equal(t, testKey, *cfg.Key())
}

func getTestConfig() Config {
	return GetTestConfig(
		WithDataBaseDSN(&testDataBaseDSN),
		WithFileStoragePath(testTempFileName),
		WithKey(&testKey),
		WithPollInterval(time.Minute),
		WithReportInterval(time.Hour),
		WithRestore(true),
		WithServerAddress(testServerAddress),
		WithStoreInterval(24*time.Hour),
	)
}

func getTestEnvironments(t *testing.T) Config {

	onceEnv = new(sync.Once)

	t.Setenv("ADDRESS", testServerAddress)
	t.Setenv("DATABASE_DSN", testDataBaseDSN)
	t.Setenv("FILE_STORAGE_PATH", testTempFileName)
	t.Setenv("KEY", testKey)
	t.Setenv("RESTORE", "true")
	t.Setenv("STORE_INTERVAL", "1")

	return GetServerConfig()
}

func getTestFlags() Config {
	onceFlags = new(sync.Once)
	initAgentFlags()
	return GetTestConfig(
		WithDataBaseDSN(&testDataBaseDSN),
		WithConfigFileName(testConfigFileName),
		WithCryptoKey(testCryptoKey),
		WithFileStoragePath(testTempFileName),
		WithKey(&testKey),
		WithPollInterval(time.Minute),
		WithReportInterval(time.Hour),
		WithRestore(true),
		WithServerAddress(testServerAddress),
		WithStoreInterval(24*time.Hour),
	)
}

func init() {
	port := 65500 + rand.Intn(34)
	testConfigFileName = fmt.Sprintf("%s/test_%018d.json", os.TempDir(), rand.Uint32())
	testCryptoKey = strings.Split("ed2b6e7e25a4c3ea1625265473a9ad19cc3a4b479137daae7895aec6202f2699dfa8556d887f0d0dac8e31f9c3ab39d84b0c7d84f7281b664cce5b2"+
		"cb4af52ea41d4c4886873989678c118e43346a854a8bcf63182b243aca15a8d75a0aefd1556d97077c0c1c6be72edc58aa2c4b946ff68e95d500b88"+
		"7a4148aedcbbfef27b3e663a54ad8b8ca7c40ce6ff3c6a63cff3ce1b34f0a190f95a31f0bf87657cc9a8fb734c15a2dc91a238e90baf8af632e07f7"+
		"ea529a9bd7660b3d3bfdbb848d55e0e95f64b55a3d1cc76ef0657ea6db1e6e6accb8f23c7ae6ac664294ca7add8332af19b85646fb3fe9257e8d028"+
		"87bfa75cb8584d4704c80e3273303d88286e8506fecddd1c23703ba361b8d0134f40492e1b1669129002d3ca7d8a2ba2623d767fe64d33599080b76"+
		"89ef8b518e73ee4b2362ac2d6e4fac5662ef32ae94a0379af07cfe5338f0d6943f8ab6220b9c84c9e369af90a344af3e4afb2a414d8ddd58e5a75ed"+
		"799dac341bdf8ca2b62b97971701165b6b8a0dda66c2211dc1ce901993be49b250f1bfba866d4e43881ebf950923b0790acce0fb3c7cf157139a4d3"+
		"e10c1068965c12fe11df4fe49dc5db61fdd2e36fc8d66f3d6757590d627818f48b040e100099eaf9c7e34f28592db5e4baa05ad488ecb7d01cbb460"+
		"e7effcd9411161b2da6008df369b5a919363f7debab33186787b181502226c657f7c4161ae69fcb4afecfa3d0e6631a8d04dfbb6bfc5f489acd37f2"+
		"c25fd15751b83dd08c26bfd391762b1340320fdb925ecd44371cd45ffb8093b3624a0a6e94345805c107ca623463e14f0580d22b86fa41bdfb8dee4"+
		"c516c5778568fcfd2d7fc8005e71783d6fc7c840cf5f0d4ca73251db14a1d500929ba28415b67133eba3cc0fe1567ca4b89c53f5eee852f029e6ba6"+
		"142373add1217e541cecceb31947e9692d8306e462646c5b52bc530b6eefc87dabdea6b9d0fe973fa2b3bacfa869349cf46f361a5eadef4a93b1593"+
		"986565b78807bbf492132f8a0cfe0789b66e96da11911411c614af48b1fced9409d32f5f9a641632e44f8efbeeeb9651607e1707b1be798e741b27d"+
		"dd6d1cbe9ffc719d42aad250274ae749362d7df88224f2d55b25e0a0362c4b80c64cee57490edcbd9bb47c66de343922bd4851ce832d6ba891dfc6c"+
		"b92ec7e08566fe0dd971fbbd8a2f9db4444c0c36576999489099da87b49c4889d78d531459868a47ebc1fb24f142175f7562b3dfb0f173cf67edd56"+
		"b3ae557fdfbf2edb29fd61e03f771a18ef41613b9f0d7f09a99ed0b9e21dcf7faca995daa968c8cfe577728ffaee917b4f00aac28ddfbd67218e8f0"+
		"5da75f1a18b246514873e753cd964daf4ec566de36510f4ade4b6da2aa1f3678b0cd9b337000e11fadbebf432ddfad05fb43bae4c267937a4571476"+
		"27103ea3c26ad89fe3a21d15442593d95740683365d06d4219fa5c8a32af6ba13a034dbde0cbfa1ec09621968f28098b70cfae9c756fbcee9418b0e"+
		"009e92eb18f5c5196494999e011a97e8f8909aefeef19116cf0a7403873865eb1d0a47d76b5117a65165d7ce6fb329ef621224c7486243385155a54"+
		"f5db0055bbd1c389549a9dd53ece44965c3bfd3558122c040f4a2a36eae8079b085d766ea9fa8b60363569dc2e363c03730be1452630fa973413241"+
		"f6ae89835a407344afcb243a7aa5efb010d74df00b1d14468eb8f7a53c2554b69b80c870ed28825892100a5c04cbf11aad703377ed2ab6eeee37751"+
		"4fcd14d600caff3c3992d30a06366aa0a333e2622575ae209969511e74be9-"+
		"57f8952a7b620b3d67eb6db54395bb5c6751741cdeb39e0cbb33a2396134bc17b54b262cd3c5749d0ed54b9a3e46e3d4e3b985516dafd5c1f359539"+
		"e5882f3b96aa3b62676bff1e8567021807bb888b5a64538f89266f70d30687b321760f1051a4adb19c048e024425309c3ebae169f4c1ba8fbf0b5e6"+
		"709a88903bd60bb1fb9ee0a26048630bd02abb326f98044a7397166b270ede8efcbec4e96cd197065b72af9a4bad76d463bc8d697be9281d3c2fe37"+
		"aa15ebbd447f5849e557af577269595c9702aef602d70aee8e82c720b3042da6b5bb2036e8cb6945210b39b3ab387ff34d120006a229499c9bbd482"+
		"442005bf31ada7346184a727756c7e5663d0e2ccfe7473c07d30a31e6e1ff426456e2df6c91b5dec01b8dac03022b11eb15be4f979fe46bec107cd8"+
		"023ff5285233943f10a807ff6dfc658b417542125994b82b636aa9974af6a25128e495bc391c5c70da4698558909f00aa61b55c815313ada3f8cfe7"+
		"3f060a1ae5ffdc957a06df5aaeb0bb8f3b6b60ca12f58cfd11a741f20ac83963c7fc73c687184240705c1f69f392e6616f6d21426f6454cc9cbbce4"+
		"ce84000e1b3f99d7dddf3c2e2b8b795d769de3b8d201200fe1a617a50c5249c7634d24be30cae6cab1a711f48067a1216f4c0d84ea290be148d6fa1"+
		"c646513e709136f0f5694bc127eb2603fb6cc8433873413a84b3fbf1c6f1e7cd1a2587c6e8bfed79440c1da1ca1683f2ff18fdf6025a35e35f54afc"+
		"63b6d0da97b8444d77bee3c6611d38cb1233bf880d59d6bfa7c5acccd16816e121e31ecd777196f8c369d55ae9dde673be31d3434bf4c5afe77ff4d"+
		"5201a94b3514054b0fb667042403e8bab91c9c9dd58d874c587424c2a37b879b87c5ecd91b88a9a2bd93eb1c341d575ae2efc32f247245d3c51df5d"+
		"c23f3e3bfb6b4ed60c284e55bf56def8cb4de7e90b0df0ec1b65510ecff54a34fa30d3e856627d6813e52899080453c4eb0dd08f511fc7c4a616fab"+
		"9e7c7154887ca76307567cb6a86e1f457640b52fa12b232f45d977423157b7caec22a80626b636eb4506f2d03b4513b542c6790dedea8c5436f3550"+
		"c5e7cb9dd1d1058c1bb6f5aaa6ed12b408e00102308692cace29efb0cc270cc79dc5946bd8c6a8f3758ab696f0bc591aaa8e3f7521d4b67e94b6788"+
		"93d92b79bdf9ff9864065becffd67b6f830b5aa3f0565fca7f8683e7a3011b7536870353c19ce5c0530523b46e46fc79b14de3091a6ec45a3f25594"+
		"9926454a3dd2add77d8824c1c2f1a03a9f477a225f89ad23d8d4ddede61563ebd7fecb278eff070dadc81d868c2f382753a9a99ad21a14f17939bf7"+
		"190c771ab15b3642e07ed158909cb48a68bd272d085d7bb14b5519893ab8b4af2a1f066a5f0739e4521a2528b923e29af4a473a44aad566d1c9e55d"+
		"1da5cb5d4146623440beb9d148f9b89220b5ced3aae050b93a53f4fe02641a105e7da7b5145cc2ff8b87054c5fcb7052b29be0ee641ad49e450cf2f"+
		"a87cd105f93d51cb994a2c4a893d5ff5ac1680a58ec1a8cb38211734071497dfc5fa8309f08e51e63d422cc098697faa856e6ddbb2dd3b34a92ed8f"+
		"98cdd6ee2ca804c266ea44931be938cb5c475afeb2e67e1fe363a63cdd4783271d43cf55bf2729e98729cd2b1326384658e44b0378a86b852ab5855"+
		"f2d7078e215ce9fac223ee34b97c0014a15aefae20fae069567086aef682f12c61a968542cfcb6f560ddb8fe97ad5826441b5513eee5afd319f2b68"+
		"48b72e65c782a6bcd2258bbc7822f68c07756b1243b26aac269ccf92c6e91", "-")
	testKey = fmt.Sprintf("%018d", rand.Uint32())
	testServerAddress = fmt.Sprintf("localhost:%d", port)
	testTempFileName = fmt.Sprintf("%s/test_%018d.txt", os.TempDir(), rand.Uint32())
}
