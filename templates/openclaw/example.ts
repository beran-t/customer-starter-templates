import { Sandbox } from 'e2b';

const tag = process.env.E2B_TEMPLATE_TAG;
const templateRef = tag ? `openclaw:${tag}` : 'openclaw';

const sbx = await Sandbox.create(templateRef, { timeoutMs: 60_000 });
try {
  const openclaw = await sbx.commands.run('openclaw --version');
  if (openclaw.exitCode !== 0) throw new Error(`openclaw check failed: ${openclaw.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}
