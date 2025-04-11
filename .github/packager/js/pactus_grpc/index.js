const modules = {};

try {
  const fs = require('fs');
  const path = require('path');
  
  fs.readdirSync(__dirname)
    .filter(file => file.endsWith('.js') && file !== 'index.js')
    .forEach(file => {
      const moduleName = path.basename(file, '.js');
      modules[moduleName] = require(`./${file}`);
    });
  
  fs.readdirSync(__dirname)
    .filter(dir => {
      const dirPath = path.join(__dirname, dir);
      return fs.existsSync(dirPath) && fs.statSync(dirPath).isDirectory();
    })
    .forEach(dir => {
      modules[dir] = {};
      fs.readdirSync(path.join(__dirname, dir))
        .filter(file => file.endsWith('.js'))
        .forEach(file => {
          const moduleName = path.basename(file, '.js');
          modules[dir][moduleName] = require(`./${dir}/${file}`);
        });
    });
} catch (err) {
  console.error('Error loading gRPC modules:', err);
}

module.exports = modules; 