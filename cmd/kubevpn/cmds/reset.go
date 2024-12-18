package cmds

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
	"k8s.io/utils/ptr"

	"github.com/wencaiwulue/kubevpn/v2/pkg/daemon"
	"github.com/wencaiwulue/kubevpn/v2/pkg/daemon/rpc"
	pkgssh "github.com/wencaiwulue/kubevpn/v2/pkg/ssh"
	"github.com/wencaiwulue/kubevpn/v2/pkg/util"
)

func CmdReset(f cmdutil.Factory) *cobra.Command {
	var sshConf = &pkgssh.SshConfig{}
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset all resource create by kubevpn in k8s cluster",
		Long: templates.LongDesc(i18n.T(`
		Reset all resource create by kubevpn in k8s cluster
		
		Reset will delete all resources create by kubevpn in k8s cluster, like deployment, service, serviceAccount...
		and it will also delete local develop docker containers, docker networks. delete hosts entry added by kubevpn,
		cleanup DNS settings.
		`)),
		Example: templates.Examples(i18n.T(`
        # Reset default namespace
		  kubevpn reset

		# Reset another namespace test
		  kubevpn reset -n test

		# Reset cluster api-server behind of bastion host or ssh jump host
		kubevpn reset --ssh-addr 192.168.1.100:22 --ssh-username root --ssh-keyfile ~/.ssh/ssh.pem

		# It also supports ProxyJump, like
		┌──────┐     ┌──────┐     ┌──────┐     ┌──────┐                 ┌────────────┐
		│  pc  ├────►│ ssh1 ├────►│ ssh2 ├────►│ ssh3 ├─────►... ─────► │ api-server │
		└──────┘     └──────┘     └──────┘     └──────┘                 └────────────┘
		kubevpn reset --ssh-alias <alias>

		# Support ssh auth GSSAPI
        kubevpn reset --ssh-addr <HOST:PORT> --ssh-username <USERNAME> --gssapi-keytab /path/to/keytab
        kubevpn reset --ssh-addr <HOST:PORT> --ssh-username <USERNAME> --gssapi-cache /path/to/cache
        kubevpn reset --ssh-addr <HOST:PORT> --ssh-username <USERNAME> --gssapi-password <PASSWORD>
		`)),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			util.InitLoggerForClient(false)
			return daemon.StartupDaemon(cmd.Context())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			bytes, ns, err := util.ConvertToKubeConfigBytes(f)
			if err != nil {
				return err
			}
			cli := daemon.GetClient(false)
			disconnect, err := cli.Disconnect(cmd.Context(), &rpc.DisconnectRequest{
				KubeconfigBytes: ptr.To(string(bytes)),
				Namespace:       ptr.To(ns),
				SshJump:         sshConf.ToRPC(),
			})
			if err != nil {
				log.Warnf("Failed to disconnect from cluter: %v", err)
			} else {
				_ = util.PrintGRPCStream[rpc.DisconnectResponse](disconnect)
			}

			req := &rpc.ResetRequest{
				KubeconfigBytes: string(bytes),
				Namespace:       ns,
				SshJump:         sshConf.ToRPC(),
			}
			resp, err := cli.Reset(cmd.Context(), req)
			if err != nil {
				return err
			}
			err = util.PrintGRPCStream[rpc.ResetResponse](resp)
			if err != nil {
				if status.Code(err) == codes.Canceled {
					return nil
				}
				return err
			}
			return nil
		},
	}

	pkgssh.AddSshFlags(cmd.Flags(), sshConf)
	return cmd
}
