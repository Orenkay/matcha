const path = require('path')
const express = require('express')
const app = express()

app.use('/dist', express.static(path.join(__dirname, '/dist')))

app.use('*', (req, res) => {
  res.sendFile(path.join(__dirname + '/index.html'))
})

app.listen(80)
