// Package portpolicy provides rule-based allow/deny evaluation for observed
// network ports.
//
// Usage:
//
//	p := portpolicy.New()
//	p.Add(portpolicy.Rule{Name: "block-telnet", Port: 23, Protocol: "tcp", Action: portpolicy.Deny})
//
//	violations := p.Check(snap)
//	for _, v := range violations {
//		log.Println(v)
//	}
//
// An Enforcer can run the check on a schedule:
//
//	e := portpolicy.NewEnforcer(p, snapFn, 30*time.Second, nil)
//	e.Run(ctx)
package portpolicy
