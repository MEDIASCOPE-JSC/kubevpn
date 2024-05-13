package cmds

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/wencaiwulue/kubevpn/v2/pkg/config"
	"github.com/wencaiwulue/kubevpn/v2/pkg/daemon"
	"github.com/wencaiwulue/kubevpn/v2/pkg/util"
)

func CmdDaemon(_ cmdutil.Factory) *cobra.Command {
	var opt = &daemon.SvrOption{}
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: i18n.T("Startup kubevpn daemon server"),
		Long:  templates.LongDesc(i18n.T(`Startup kubevpn daemon server`)),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			b := make([]byte, 32)
			if _, err := rand.Read(b); err != nil {
				return err
			}
			opt.ID = base64.URLEncoding.EncodeToString(b)

			if opt.IsSudo {
				go util.StartupPProf(config.SudoPProfPort)
			} else {
				go util.StartupPProf(config.PProfPort)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer opt.Stop()
			defer func() {
				if errors.Is(err, http.ErrServerClosed) {
					err = nil
				}
			}()
			return opt.Start(cmd.Context())
		},
		Hidden:                true,
		DisableFlagsInUseLine: true,
	}
	cmd.Flags().BoolVar(&opt.IsSudo, "sudo", false, "is sudo or not")
	return cmd
}
