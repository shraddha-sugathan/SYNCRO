import { copyFileSync } from 'node:fs'
copyFileSync('dist/app.html', 'index.html')
console.log('✔ frontend/index.html updated (served by main.go at "/")')
