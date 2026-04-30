var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

document.addEventListener('DOMContentLoaded', () => {
    const btn = document.querySelector('#theme-toggle-btn');
    const btnText = btn.querySelector('.text');
    const btnIcon = btn.querySelector('.icon');
    const body = document.body;

    // Function to update button UI
    const updateBtnUI = (isDark) => {
        btnText.textContent = isDark ? 'Light Mode' : 'Dark Mode';
        btnIcon.textContent = isDark ? '☀️' : '🌙';
    };

    // Check storage on load
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark-mode') {
        body.classList.add('dark-mode');
        updateBtnUI(true);
    }

    btn.addEventListener('click', () => {
        const isDarkMode = body.classList.toggle('dark-mode');
        
        // Save preference
        localStorage.setItem('theme', isDarkMode ? 'dark-mode' : 'light-mode');
        
        // Update Button
        updateBtnUI(isDarkMode);
    });
});

// document.addEventListener('DOMContentLoaded', () => {
//     const toggleSwitch = document.querySelector('#theme-toggle');
//     const body = document.body;

//     // 1. Check local storage on load
//     const currentTheme = localStorage.getItem('theme');
//     if (currentTheme === 'dark-mode') {
//         body.classList.add('dark-mode');
//         toggleSwitch.checked = true;
//     }

//     // 2. Listen for changes
//     toggleSwitch.addEventListener('change', () => {
//         if (toggleSwitch.checked) {
//             body.classList.add('dark-mode');
//             localStorage.setItem('theme', 'dark-mode');
//         } else {
//             body.classList.remove('dark-mode');
//             localStorage.setItem('theme', 'light-mode');
//         }
//     });
// });