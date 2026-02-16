import { Sandbox } from 'e2b';

const tag = process.env.E2B_TEMPLATE_TAG;
const templateRef = tag ? `amp-code:${tag}` : 'amp-code';

const sbx = await Sandbox.create(templateRef, { timeoutMs: 60_000 });
try {
  const amp = await sbx.commands.run('amp --version');
  if (amp.exitCode !== 0) throw new Error(`amp check failed: ${amp.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}
