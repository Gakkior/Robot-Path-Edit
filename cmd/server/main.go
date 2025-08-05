// Package main 鏄簲鐢ㄧ▼搴忕殑鍏ュ彛鐐?
//
// 鏋舵瀯璁捐鍙傝€冿細
// - Kubernetes API Server鐨勫惎鍔ㄦā寮?
// - Grafana Server鐨勯厤缃鐞?
// - Docker鐨勫懡浠よ鎺ュ彛璁捐
//
// 璁捐鐞嗗康锛?
// 1. 鍗曚竴鑱岃矗锛歮ain鍑芥暟鍙礋璐ｅ惎鍔紝涓嶅寘鍚笟鍔￠€昏緫
// 2. 閰嶇疆椹卞姩锛氶€氳繃閰嶇疆鏂囦欢鍜岀幆澧冨彉閲忔帶鍒惰涓?
// 3. 浼橀泤鍚仠锛氭敮鎸佷俊鍙峰鐞嗗拰璧勬簮娓呯悊
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"robot-path-editor/internal/app"
	"robot-path-editor/internal/config"
	"robot-path-editor/pkg/logger"
)

var (
	// 鐗堟湰淇℃伅锛屾瀯寤烘椂娉ㄥ叆
	version   = "dev"
	buildTime = "unknown"
	gitHash   = "unknown"
)

func main() {
	// 鍒涘缓鏍瑰懡浠?- 鍙傝€僰ubectl鐨勫懡浠ょ粨鏋?
	rootCmd := &cobra.Command{
		Use:     "robot-path-editor",
		Short:   "鏈哄櫒浜鸿矾寰勭紪杈戝櫒 - 閫氱敤鐨勭偣浣嶅拰璺緞绠＄悊宸ュ叿",
		Long:    `涓€涓幇浠ｅ寲鐨勪笁绔吋瀹规満鍣ㄤ汉璺緞缂栬緫鍣紝鏀寔鍙鍖栫紪杈戝拰鏁版嵁搴撶鐞嗐€俙,
		Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, gitHash),
		RunE:    runServer,
	}

	// 娣诲姞鍛戒护琛屽弬鏁?- 鍙傝€僄rafana鐨勯厤缃€夐」
	rootCmd.PersistentFlags().String("config", "", "閰嶇疆鏂囦欢璺緞")
	rootCmd.PersistentFlags().String("log-level", "info", "鏃ュ織绾у埆 (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("addr", ":8080", "鏈嶅姟鐩戝惉鍦板潃")
	rootCmd.PersistentFlags().String("db-type", "sqlite", "鏁版嵁搴撶被鍨?(sqlite, mysql)")
	rootCmd.PersistentFlags().String("db-dsn", "./data/app.db", "鏁版嵁搴撹繛鎺ュ瓧绗︿覆")

	// 缁戝畾鍙傛暟鍒皏iper - 鍙傝€僈ubernetes鐨勯厤缃鐞?
	viper.BindPFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("搴旂敤鍚姩澶辫触")
	}
}

// runServer 鍚姩鏈嶅姟鍣?
// 鍙傝€僄rafana Server鐨勫惎鍔ㄦ祦绋?
func runServer(cmd *cobra.Command, args []string) error {
	// 1. 鍔犺浇閰嶇疆
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("鍔犺浇閰嶇疆澶辫触: %w", err)
	}

	// 2. 鍒濆鍖栨棩蹇?- 鍙傝€僈ubernetes鐨勭粨鏋勫寲鏃ュ織
	logger.Init(cfg.Logger)

	log := logrus.WithFields(logrus.Fields{
		"version":    version,
		"build_time": buildTime,
		"git_hash":   gitHash,
	})

	log.Info("鏈哄櫒浜鸿矾寰勭紪杈戝櫒鍚姩涓?..")

	// 3. 鍒涘缓搴旂敤瀹炰緥 - 渚濊禆娉ㄥ叆妯″紡锛屽弬鑰僓ber FX
	application, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("鍒涘缓搴旂敤瀹炰緥澶辫触: %w", err)
	}

	// 4. 璁剧疆淇″彿澶勭悊 - 鍙傝€僁ocker鐨勪紭闆呭仠鏈?
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 鐩戝惉绯荤粺淇″彿
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 5. 鍚姩搴旂敤
	go func() {
		if err := application.Start(ctx); err != nil {
			log.WithError(err).Error("搴旂敤鍚姩澶辫触")
			cancel()
		}
	}()

	log.WithField("addr", cfg.Server.Addr).Info("鏈嶅姟鍣ㄥ惎鍔ㄦ垚鍔?)

	// 6. 绛夊緟鍋滄淇″彿
	select {
	case sig := <-sigChan:
		log.WithField("signal", sig).Info("鏀跺埌鍋滄淇″彿锛屽紑濮嬩紭闆呭叧闂?..")
	case <-ctx.Done():
		log.Info("搴旂敤涓婁笅鏂囪鍙栨秷")
	}

	// 7. 浼橀泤鍏抽棴 - 鍙傝€僈ubernetes鐨勪紭闆呯粓姝?
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := application.Stop(shutdownCtx); err != nil {
		log.WithError(err).Error("搴旂敤鍏抽棴鏃跺彂鐢熼敊璇?)
		return err
	}

	log.Info("搴旂敤宸叉垚鍔熷叧闂?)
	return nil
}
