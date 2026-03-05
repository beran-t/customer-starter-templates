import { Sandbox } from 'e2b';

const tag = process.env.E2B_TEMPLATE_TAG;
const templateRef = tag ? `tbench:${tag}` : 'tbench';

const sbx = await Sandbox.create(templateRef, { timeoutMs: 60_000 });
try {
  const harbor = await sbx.commands.run('harbor --version');
  if (harbor.exitCode !== 0) throw new Error(`harbor check failed: ${harbor.stderr}`);

  const tb = await sbx.commands.run('tb run --help');
  if (tb.exitCode !== 0) throw new Error(`terminal-bench check failed: ${tb.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}
