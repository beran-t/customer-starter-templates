import { Sandbox } from 'e2b';

const tag = process.env.E2B_TEMPLATE_TAG;
const templateRef = tag ? `codex:${tag}` : 'codex';

const sbx = await Sandbox.create(templateRef, { timeoutMs: 60_000 });
try {
  const codex = await sbx.commands.run('codex --version');
  if (codex.exitCode !== 0) throw new Error(`codex check failed: ${codex.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}
