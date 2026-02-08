import { Sandbox } from 'e2b';

const sbx = await Sandbox.create('e2b-tbench', { timeoutMs: 60_000 });
try {
  // Verify tmux is installed
  const tmux = await sbx.commands.run('tmux -V');
  if (tmux.exitCode !== 0) throw new Error(`tmux check failed: ${tmux.stderr}`);

  // Verify git is installed
  const git = await sbx.commands.run('git --version');
  if (git.exitCode !== 0) throw new Error(`git check failed: ${git.stderr}`);

  // Verify Docker is available
  const docker = await sbx.commands.run('docker --version');
  if (docker.exitCode !== 0) throw new Error(`docker check failed: ${docker.stderr}`);

  // Verify uv is installed
  const uv = await sbx.commands.run('uv --version');
  if (uv.exitCode !== 0) throw new Error(`uv check failed: ${uv.stderr}`);

  // Verify harbor is installed
  const harbor = await sbx.commands.run('harbor --version');
  if (harbor.exitCode !== 0) throw new Error(`harbor check failed: ${harbor.stderr}`);

  // Verify terminal-bench is installed
  const tb = await sbx.commands.run('tb run --help');
  if (tb.exitCode !== 0) throw new Error(`terminal-bench check failed: ${tb.stderr}`);

  console.log('All checks passed.');
} finally {
  await sbx.kill();
}
