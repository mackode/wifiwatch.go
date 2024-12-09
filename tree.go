package main

import (
  "fmt"
  "github.com/gdamore/tcell/v2"
  "github.com/rivo/tview"
)

func tree() {
  db, err := NewDB()
  if err != nil {
    panic(err)
  }

  root := tview.NewTreeNode("Wifiwatch v1.0").SetColor(tcell.ColorRed)
  tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
  oldMAC := ""
  var node *tview.TreeNode

  probes, err := db.Probes()
  if err != nil {
    panic(err)
  }

  for _, p := range probes {
    if node == nil || oldMAC != p.Device.MAC {
      node = tview.NewTreeNode(fmt.Sprintf("%s %s", p.Device.Product, p.Device.MAC))
      root.AddChild(node)
      node.SetColor(tcell.ColorGreen)
    }

    n := tview.NewTreeNode(fmt.Sprintf("%s %s", p.IP, p.Timestamp))
    node.AddChild(n)
    oldMAC = p.Device.MAC
  }

  err = tview.NewApplication().SetRoot(tree, true).Run()
  if err != nil {
    panic(err)
  }
}
