package cmds

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestName111(t *testing.T) {
	var processBar = []string{
		"00%: [                                          ]",
		"05%: [##                                        ]",
		"10%: [####                                      ]",
		"15%: [######                                    ]",
		"20%: [########                                  ]",
		"25%: [##########                                ]",
		"30%: [############                              ]",
		"35%: [##############                            ]",
		"40%: [################                          ]",
		"45%: [##################                        ]",
		"50%: [####################                      ]",
		"55%: [######################                    ]",
		"60%: [########################                  ]",
		"65%: [##########################                ]",
		"70%: [############################              ]",
		"75%: [##############################            ]",
		"80%: [################################          ]",
		"85%: [##################################        ]",
		"90%: [####################################      ]",
		"95%: [######################################    ]",
		"100%:[##########################################]\n",
	}
	for idx, val := range processBar {
		fmt.Printf("[%d:1H:2K] \r \a%s", idx, val)
		time.Sleep(1 * time.Millisecond * 200)
	}
}

func TestAAAAA(t *testing.T) {
	writer := io.MultiWriter(os.Stdout, os.Stderr)
	for i := 0; i < 10; i++ {
		fmt.Fprintf(writer, "\r[%d:1H:2K]", i)
		time.Sleep(1 * time.Millisecond * 200)
	}
}
