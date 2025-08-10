document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('shutdown-modal');
    const yesBtn = document.getElementById('shutdown-yes-btn');
    const noBtn = document.getElementById('shutdown-no-btn');
    const shutdownBtn = document.getElementById('shutdown-btn');

    if (shutdownBtn) {
        shutdownBtn.addEventListener('click', () => {
            if(modal) modal.style.display = 'flex';
        });
    }

    if (noBtn) {
        noBtn.addEventListener('click', () => {
            if(modal) modal.style.display = 'none';
        });
    }

    if (modal) {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) { // Click on overlay
                modal.style.display = 'none';
            }
        });
    }

    if (yesBtn) {
        yesBtn.addEventListener('click', async () => {
            try {
                await fetch('/api/shutdown', { method: 'POST' });
                document.body.innerHTML = '<div style="padding: 20px; text-align: center;">Server has been shut down. You can close this window.</div>';
            } catch (error) {
                if(modal) {
                    const p = modal.querySelector('p');
                    if (p) {
                        p.textContent = 'Failed to send shutdown signal. Please close the terminal manually.';
                        p.style.color = 'red';
                    }
                }
            }
        });
    }
});
