/*
	cmd 控制 cli app的接口
*/
package cmd

type Control interface {
	Exit()
	ClearHistory()
}
