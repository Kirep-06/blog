(function () {
  const api = window.blogApi;

  async function loadFilters() {
    const [cats, tags] = await Promise.all([api.listCategories(), api.listTags()]);
    const catSel = document.getElementById('category');
    const tagSel = document.getElementById('tag');
    for (const c of cats.data || []) {
      const op = document.createElement('option'); op.value = c.slug; op.textContent = c.name; catSel.appendChild(op);
    }
    for (const t of tags.data || []) {
      const op = document.createElement('option'); op.value = t.slug; op.textContent = t.name; tagSel.appendChild(op);
    }
  }

  function card(post) {
    const el = document.createElement('article');
    el.className = 'card';
    el.innerHTML = `
      <h3><a href="/frontend/post.html?slug=${encodeURIComponent(post.slug)}">${post.title}</a></h3>
      ${post.cover_url ? `<img class="post-cover" src="${post.cover_url}" alt="cover" />` : ''}
      <p>${post.summary || ''}</p>
      <div class="small">${new Date(post.created_at).toLocaleString()} · ${post.category?.name || '未分类'}</div>
      <div>${(post.tags || []).map(t => `<span class="badge">#${t.name}</span>`).join('')}</div>`;
    return el;
  }

  async function loadPosts() {
    const errorEl = document.getElementById('error');
    errorEl.textContent = '';
    const params = {
      page: document.getElementById('page').value || '1',
      page_size: document.getElementById('page_size').value || '10',
      q: document.getElementById('q').value,
      category: document.getElementById('category').value,
      tag: document.getElementById('tag').value,
    };
    const qs = new URLSearchParams(params);
    history.replaceState(null, '', `?${qs.toString()}`);

    try {
      const res = await api.listPosts(params);
      const list = document.getElementById('posts');
      list.innerHTML = '';
      (res.data || []).forEach(p => list.appendChild(card(p)));
      document.getElementById('total').textContent = `共 ${res.total || 0} 篇`;
    } catch (e) {
      errorEl.textContent = e.message;
    }
  }

  async function loadPostDetail() {
    const content = document.getElementById('post-detail');
    if (!content) return;
    const slug = new URLSearchParams(location.search).get('slug');
    if (!slug) return (content.textContent = '缺少 slug 参数');
    try {
      const { data: post } = await api.getPost(slug);
      document.title = post.title;
      content.innerHTML = `
      <h1>${post.title}</h1>
      <div class="small">${new Date(post.created_at).toLocaleString()} · ${post.category?.name || '未分类'}</div>
      ${post.cover_url ? `<img class="post-cover" src="${post.cover_url}" alt="cover"/>` : ''}
      <div>${(post.tags || []).map(t => `<span class="badge">#${t.name}</span>`).join('')}</div>
      <article class="card">${post.content_html || ''}</article>`;
    } catch (e) {
      content.innerHTML = `<p id="error">${e.message}</p>`;
    }
  }

  async function bootstrapListPage() {
    if (!document.getElementById('posts')) return;
    const query = new URLSearchParams(location.search);
    ['page', 'page_size', 'q', 'category', 'tag'].forEach((k) => {
      const el = document.getElementById(k);
      if (el && query.get(k)) el.value = query.get(k);
    });

    await loadFilters();
    await loadPosts();
    document.getElementById('search-btn').addEventListener('click', loadPosts);
  }

  bootstrapListPage();
  loadPostDetail();
})();
