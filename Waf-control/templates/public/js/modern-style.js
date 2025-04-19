
// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 添加导航栏激活状态
    highlightActiveNavItem();
    
    // 添加表格行悬停效果
    addTableRowHoverEffect();
    
    // 添加卡片悬停效果
    addCardHoverEffect();
    
    // 添加按钮点击波纹效果
    addButtonRippleEffect();
    
    // 添加页面过渡动画
    addPageTransitionEffect();
    
    // 添加表单输入焦点效果
    addFormFocusEffect();
    
    // 添加警告框自动消失效果
    setupAlertAutoDismiss();
});

/**
 * 高亮当前激活的导航项
 */
function highlightActiveNavItem() {
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.navbar-nav .nav-link');
    
    navLinks.forEach(link => {
        const href = link.getAttribute('href');
        if (currentPath.includes(href) && href !== '/logout') {
            link.classList.add('active');
            link.style.backgroundColor = 'var(--hover-color)';
            link.style.opacity = '1';
        }
    });
}

/**
 * 添加表格行悬停效果
 */
function addTableRowHoverEffect() {
    const tableRows = document.querySelectorAll('tbody tr');
    
    tableRows.forEach(row => {
        row.addEventListener('mouseenter', function() {
            this.style.backgroundColor = 'var(--hover-color)';
            this.style.transform = 'scale(1.01)';
            this.style.transition = 'all 0.2s ease';
        });
        
        row.addEventListener('mouseleave', function() {
            this.style.backgroundColor = '';
            this.style.transform = 'scale(1)';
        });
    });
}

/**
 * 添加卡片悬停效果
 */
function addCardHoverEffect() {
    const cards = document.querySelectorAll('.card');
    
    cards.forEach(card => {
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-8px)';
            this.style.boxShadow = '0 12px 30px rgba(0, 0, 0, 0.4)';
        });
        
        card.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
            this.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.2)';
        });
    });
}

/**
 * 添加按钮点击波纹效果
 */
function addButtonRippleEffect() {
    const buttons = document.querySelectorAll('.btn');
    
    buttons.forEach(button => {
        button.addEventListener('click', function(e) {
            const rect = this.getBoundingClientRect();
            const x = e.clientX - rect.left;
            const y = e.clientY - rect.top;
            
            const ripple = document.createElement('span');
            ripple.classList.add('ripple');
            ripple.style.left = `${x}px`;
            ripple.style.top = `${y}px`;
            
            this.appendChild(ripple);
            
            setTimeout(() => {
                ripple.remove();
            }, 600);
        });
    });
}

/**
 * 添加页面过渡动画
 */
function addPageTransitionEffect() {
    const contentElements = document.querySelectorAll('.container > *:not(nav)');
    
    contentElements.forEach((element, index) => {
        element.style.opacity = '0';
        element.style.transform = 'translateY(20px)';
        element.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
        
        setTimeout(() => {
            element.style.opacity = '1';
            element.style.transform = 'translateY(0)';
        }, 100 + (index * 50));
    });
}

/**
 * 添加表单输入焦点效果
 */
function addFormFocusEffect() {
    const formControls = document.querySelectorAll('.form-control');
    
    formControls.forEach(input => {
        // 创建标签动画效果
        const formGroup = input.closest('.form-group');
        if (formGroup) {
            const label = formGroup.querySelector('label');
            if (label) {
                label.style.transition = 'all 0.3s ease';
                
                input.addEventListener('focus', function() {
                    label.style.color = 'var(--accent-color)';
                    label.style.transform = 'scale(1.05)';
                    label.style.transformOrigin = 'left';
                });
                
                input.addEventListener('blur', function() {
                    label.style.color = '';
                    label.style.transform = 'scale(1)';
                });
            }
        }
        
        // 输入框焦点效果
        input.addEventListener('focus', function() {
            this.style.transform = 'scale(1.02)';
            this.style.transformOrigin = 'left';
        });
        
        input.addEventListener('blur', function() {
            this.style.transform = 'scale(1)';
        });
    });
}

/**
 * 设置警告框自动消失
 */
function setupAlertAutoDismiss() {
    const alerts = document.querySelectorAll('.alert');
    
    alerts.forEach(alert => {
        setTimeout(() => {
            alert.style.opacity = '0';
            alert.style.transform = 'translateY(-20px)';
            alert.style.transition = 'all 0.5s ease';
            
            setTimeout(() => {
                alert.remove();
            }, 500);
        }, 5000);
    });
}