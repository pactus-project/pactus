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
} catch (err) {
  console.error('Error loading gRPC modules:', err);
}

module.exports = modules;
