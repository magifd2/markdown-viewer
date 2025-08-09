const container = document.getElementById('tree-container');

async function fetchAndRender(path, parentElement) {
    try {
        const response = await fetch('/api/list?path=' + encodeURIComponent(path));
        if (!response.ok) throw new Error('Failed to fetch');
        const items = await response.json();

        if (items === null || items.length === 0) {
            const parentIcon = parentElement.closest('.node-container').querySelector('.icon');
            if(parentIcon) parentIcon.classList.add('empty');
            return;
        }

        const ul = document.createElement('ul');
        items.forEach(item => {
            const li = document.createElement('li');
            li.classList.add('node-container');

            const label = document.createElement('div');
            label.classList.add('node-label');

            const icon = document.createElement('span');
            icon.classList.add('icon');

            const link = document.createElement('a');
            link.textContent = item.name;

            if (item.is_dir) {
                icon.textContent = '▶';
                link.href = 'javascript:void(0)';
                label.addEventListener('click', (e) => {
                    e.stopPropagation();
                    toggleNode(li);
                });
                // Pass path to li element for toggleNode to use
                li.dataset.path = item.path;
            } else {
                icon.innerHTML = '&#128196;';
                link.href = 'javascript:void(0)';
                label.addEventListener('click', (e) => {
                    e.preventDefault();
                    window.parent.document.querySelector('iframe[name="content_frame"]').src = `/view${item.path}`;
                });
            }
            
            label.appendChild(icon);
            label.appendChild(link);
            li.appendChild(label);
            ul.appendChild(li);
        });
        parentElement.appendChild(ul);
    } catch (error) {
        console.error('Error fetching directory:', error);
        parentElement.textContent = 'Error loading directory.';
    }
}

function toggleNode(li) {
    const childrenContainer = li.querySelector('.children');
    const icon = li.querySelector('.icon');
    const path = li.dataset.path;

    if (icon.classList.contains('empty')) return;

    if (childrenContainer) {
        childrenContainer.classList.toggle('expanded');
        icon.textContent = childrenContainer.classList.contains('expanded') ? '▼' : '▶';
    } else {
        const newChildrenContainer = document.createElement('div');
        newChildrenContainer.classList.add('children', 'expanded');
        li.appendChild(newChildrenContainer);
        fetchAndRender(path, newChildrenContainer);
        icon.textContent = '▼';
    }
}

// Initial load
fetchAndRender('/', container);
