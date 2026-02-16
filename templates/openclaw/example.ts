import { Sandbox } from 'e2b';
import { readFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const testScript = readFileSync(join(__dirname, 'test.sh'), 'utf-8');

const sbx = await Sandbox.create('openclaw', { timeoutMs: 60_000 });
try {
  await sbx.files.write('/tmp/test.sh', testScript);
  const result = await sbx.commands.run('bash /tmp/test.sh');
  if (result.exitCode !== 0) throw new Error(`Test failed:\n${result.stderr}`);
  console.log(result.stdout);
} finally {
  await sbx.kill();
}
