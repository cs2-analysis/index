import cp from 'child_process';
import fs from 'fs';

function run(branch, file) {
  return new Promise((resolve, reject) => {
    console.log(`Generating data for ${branch}...`);

    const writeStream = fs.createWriteStream(file);

    const proc = cp.execFile('go', ['run', 'main.go', '-b', branch], {
      stdio: ['ignore', 'pipe', 'pipe'],
      cwd: import.meta.dirname
    });

    proc.stdout.pipe(writeStream);
    proc.stderr.pipe(process.stdout);

    proc.on('exit', (code) => {
      if (code === 0) {
        resolve();
      } else {
        reject(new Error(`Process exited with code ${code}`));
      }
    });
  });
}

run('windows', 'data/windows.json')
  .then(() => run('linux', 'data/linux.json'))
  .catch(err => {
    console.error(err);
    process.exit(1);
  })
  .then(() => {
    console.log('Done');
    process.exit(0);
  });
