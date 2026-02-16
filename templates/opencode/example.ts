import { Sandbox } from 'e2b';

const tag = process.env.E2B_TEMPLATE_TAG;
const templateRef = tag ? `opencode:${tag}` : 'opencode';

const sbx = await Sandbox.create(templateRef, { timeoutMs: 60_000 });
try {
  const opencode = await sbx.commands.run('opencode --version');
  if (opencode.exitCode !== 0) throw new Error(`opencode check failed: ${opencode.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}
