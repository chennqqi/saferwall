package main

import (
	"flag"
	"log"
	"os"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	avastclient "github.com/saferwall/saferwall/pkg/grpc/multiav/avast/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/avast/proto"
	aviraclient "github.com/saferwall/saferwall/pkg/grpc/multiav/avira/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/avira/proto"
	bitdefenderclient "github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/proto"
	clamavclient "github.com/saferwall/saferwall/pkg/grpc/multiav/clamav/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/clamav/proto"
	comodoclient "github.com/saferwall/saferwall/pkg/grpc/multiav/comodo/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/comodo/proto"
	esetclient "github.com/saferwall/saferwall/pkg/grpc/multiav/eset/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/eset/proto"
	fsecureclient "github.com/saferwall/saferwall/pkg/grpc/multiav/fsecure/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/fsecure/proto"
	kasperskyclient "github.com/saferwall/saferwall/pkg/grpc/multiav/kaspersky/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/kaspersky/proto"
	mcafeeclient "github.com/saferwall/saferwall/pkg/grpc/multiav/mcafee/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/mcafee/proto"
	sophosclient "github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/proto"
	symantecclient "github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/proto"
	windefenderclient "github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"github.com/saferwall/saferwall/pkg/utils"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "172.17.0.2:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.saferwall.com", "The server name use to verify the hostname returned by TLS handshake")
	filePath           = flag.String("path", "", "The file path or directory to scan")
	engine           = flag.String("engine", "", "The antivirus engine used to scan the file")
)

// parseFlags parses the cmd line flags to create grpc conn.
func parseFlags() (string, []grpc.DialOption, string, string) {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	if *filePath == "" || *engine == "" {
		flag.Usage()
		os.Exit(0)
	}
	return *serverAddr, opts, *filePath, *engine
}

func scan(engine string, filePath string, conn *grpc.ClientConn) {

	var res multiav.ScanResult
	var err error

	switch engine {
	case "avast":
		res, err = avastclient.ScanFile(avast_api.NewAvastScannerClient(conn), filePath)
	case "avira":
		res, err = aviraclient.ScanFile(avira_api.NewAviraScannerClient(conn), filePath)
	case "bitdefender":
		res, err = bitdefenderclient.ScanFile(bitdefender_api.NewBitdefenderScannerClient(conn), filePath)
	case "clamav":
		res, err = clamavclient.ScanFile(clamav_api.NewClamAVScannerClient(conn), filePath)
	case "comodo":
		res, err = comodoclient.ScanFile(comodo_api.NewComodoScannerClient(conn), filePath)
	case "eset":
		res, err = esetclient.ScanFile(eset_api.NewEsetScannerClient(conn), filePath)
	case "fsecure":
		res, err = fsecureclient.ScanFile(fsecure_api.NewFSecureScannerClient(conn), filePath)
	case "kaspersky":
		res, err = kasperskyclient.ScanFile(kaspersky_api.NewKasperskyScannerClient(conn), filePath)
	case "mcafee":
		res, err = mcafeeclient.ScanFile(mcafee_api.NewMcAfeeScannerClient(conn), filePath)
	case "symantec":
		res, err = symantecclient.ScanFile(symantec_api.NewSymantecScannerClient(conn), filePath)
	case "sophos":
		res, err = sophosclient.ScanFile(sophos_api.NewSophosScannerClient(conn), filePath)
	case "windefender":
		res, err = windefenderclient.ScanFile(windefender_api.NewWinDefenderScannerClient(conn), filePath)
	}

	if err != nil {
		log.Printf("Failed to scan file [%s]: %v", engine, err)
	}
	
	log.Print(filePath, res)
}



func main() {
	// Parse command line args.
	serverAddr, _, filePath, engine := parseFlags()

	// Get grpc client conn.
	conn, err := multiav.GetClientConn(serverAddr)
	if err != nil {
		log.Fatalf("fail to dial for engine %s : %v", engine, err)
	}
	defer conn.Close()

	files, err := utils.WalkAllFilesInDir(filePath)
	if err != nil {
		log.Fatalf("fail to walk dir %s : %v", filePath, err)
	}

	// Scan the file.
	for _, file := range files {
		scan(engine, file, conn)
	}
}