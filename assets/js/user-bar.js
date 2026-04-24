const loginName = localStorage.getItem('loginName');

const bar = document.createElement('div');
bar.id = 'user-bar';
bar.textContent = loginName ?? '';
bar.style.cssText = 'position:fixed;top:12px;right:20px;color:#4a9eff;font-size:14px';
document.body.prepend(bar);
