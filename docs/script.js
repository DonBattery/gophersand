const spriteData = {
    "Empty": {
        "inactive_sprite": {
            "x": 0,
            "y": 0,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 24,
            "y": 0,
            "w": 24,
            "h": 24
        }
    },
    "Stone": {
        "inactive_sprite": {
            "x": 48,
            "y": 0,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 72,
            "y": 0,
            "w": 24,
            "h": 24
        }
    },
    "Wood": {
        "inactive_sprite": {
            "x": 96,
            "y": 0,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 120,
            "y": 0,
            "w": 24,
            "h": 24
        }
    },
    "Sand": {
        "inactive_sprite": {
            "x": 144,
            "y": 0,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 168,
            "y": 0,
            "w": 24,
            "h": 24
        }
    },
    "Water": {
        "inactive_sprite": {
            "x": 192,
            "y": 0,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 216,
            "y": 0,
            "w": 24,
            "h": 24
        }
    },
    "Seed": {
        "inactive_sprite": {
            "x": 0,
            "y": 24,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 24,
            "y": 24,
            "w": 24,
            "h": 24
        }
    },
    "Ant": {
        "inactive_sprite": {
            "x": 48,
            "y": 24,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 72,
            "y": 24,
            "w": 24,
            "h": 24
        }
    },
    "Acid": {
        "inactive_sprite": {
            "x": 96,
            "y": 24,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 120,
            "y": 24,
            "w": 24,
            "h": 24
        }
    },
    "Fire": {
        "inactive_sprite": {
            "x": 144,
            "y": 24,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 168,
            "y": 24,
            "w": 24,
            "h": 24
        }
    },
    "Ice": {
        "inactive_sprite": {
            "x": 192,
            "y": 24,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 216,
            "y": 24,
            "w": 24,
            "h": 24
        }
    },
    "Start": {
        "inactive_sprite": {
            "x": 0,
            "y": 48,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 24,
            "y": 48,
            "w": 24,
            "h": 24
        }
    },
    "Stop": {
        "inactive_sprite": {
            "x": 48,
            "y": 48,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 72,
            "y": 48,
            "w": 24,
            "h": 24
        }
    },
    "Erase": {
        "inactive_sprite": {
            "x": 96,
            "y": 48,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 120,
            "y": 48,
            "w": 24,
            "h": 24
        }
    },
    "Gen": {
        "inactive_sprite": {
            "x": 144,
            "y": 48,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 168,
            "y": 48,
            "w": 24,
            "h": 24
        }
    },
    "SizePlus": {
        "inactive_sprite": {
            "x": 192,
            "y": 48,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 216,
            "y": 48,
            "w": 24,
            "h": 24
        }
    },
    "SizeMinus": {
        "inactive_sprite": {
            "x": 0,
            "y": 72,
            "w": 24,
            "h": 24
        },
        "active_sprite": {
            "x": 24,
            "y": 72,
            "w": 24,
            "h": 24
        }
    },
    GitHub: {
        "inactive_sprite": {
            "x": 48,
            "y": 81,
            "w": 41,
            "h": 15
        },
        "active_sprite": {
            "x": 48,
            "y": 96,
            "w": 41,
            "h": 15
        }
    },
    Title: {
        "inactive_sprite": {
            "x": 48,
            "y": 72,
            "w": 50,
            "h": 9
        },
        "active_sprite": {
            "x": 98,
            "y": 72,
            "w": 50,
            "h": 9
        }
    }
};

// Sprite sheet dimensions in pixels
const SPRITE_SHEET_WIDTH = 240;
const SPRITE_SHEET_HEIGHT = 120;

/**
 * Sprite-backed button using regions from a shared sprite sheet.
 */
class SpriteButton {
    /**
     * @param {Object} options
     * @param {string} options.spriteSheetUrl - URL/path to the sprite sheet image.
     * @param {Object} options.inactiveSprite - { x, y, w, h } for inactive state.
     * @param {Object} options.activeSprite - { x, y, w, h } for active state.
     * @param {Function} [options.onClick] - Callback invoked on click: (event, buttonInstance).
     * @param {number} [options.scale=1] - Optional scale factor for button size.
     * @param {number} [options.sheetWidth=SPRITE_SHEET_WIDTH] - Full sprite sheet width.
     * @param {number} [options.sheetHeight=SPRITE_SHEET_HEIGHT] - Full sprite sheet height.
     */
    constructor({ spriteSheetUrl, inactiveSprite, activeSprite, onClick, scale = 1, sheetWidth = SPRITE_SHEET_WIDTH, sheetHeight = SPRITE_SHEET_HEIGHT }) {
        this.spriteSheetUrl = spriteSheetUrl;
        this.inactiveSprite = inactiveSprite;
        this.activeSprite = activeSprite;
        this.onClick = onClick;
        this.scale = scale;
        this.sheetWidth = sheetWidth;
        this.sheetHeight = sheetHeight;
        this.isActive = false;

        this.element = document.createElement('button');
        this.element.type = 'button';
        this.element.classList.add('sprite-button');

        if (typeof this.onClick === 'function') {
            this.element.addEventListener('click', (event) => {
                this.onClick(event, this);
            });
        }

        this._applySprite();
    }

    /**
     * Apply current sprite (inactive/active) to the button element.
     * Keeps the sprite's intrinsic aspect ratio, scaled by this.scale.
     * @private
     */
    _applySprite() {
        const frame = this.isActive ? this.activeSprite : this.inactiveSprite;
        if (!frame) {
            return;
        }

        const width = frame.w * this.scale;
        const height = frame.h * this.scale;

        this.element.style.width = `${width}px`;
        this.element.style.height = `${height}px`;
        this.element.style.backgroundImage = `url(${this.spriteSheetUrl})`;
        this.element.style.backgroundRepeat = 'no-repeat';
        // Scale the whole sprite sheet so the selected frame fills the button
        this.element.style.backgroundSize = `${this.sheetWidth * this.scale}px ${this.sheetHeight * this.scale}px`;
        this.element.style.backgroundPosition = `-${frame.x * this.scale}px -${frame.y * this.scale}px`;
    }

    /**
     * Set active state explicitly.
     * @param {boolean} isActive
     */
    setActive(isActive) {
        this.isActive = !!isActive;
        this._applySprite();
    }

    /**
     * Update the scale factor and re-apply sprite.
     * @param {number} scale
     */
    setScale(scale) {
        this.scale = scale;
        this._applySprite();
    }

    /**
     * Toggle between inactive and active sprites.
     */
    toggle() {
        this.setActive(!this.isActive);
    }

    /**
     * Attach the button to a parent DOM element.
     * @param {HTMLElement} parent
     * @returns {SpriteButton}
     */
    mount(parent) {
        if (parent && parent.appendChild) {
            parent.appendChild(this.element);
        }
        return this;
    }
}

// Define the App object
const App = {
    // Properties
    Width: 0,
    Height: 0,
    CurrentLayout: '',
    Buttons: [],
    FooterButtons: null,
    BrushButtons: [],
    CurrentBrushName: 'Sand',
    IsRunning: false,
    ActiveModalType: null, // 'erase' | 'gen' | null

    /**
     * Unified event emitter for communicating with the Go/WASM game.
     * For now it simply logs the event string.
     * @param {string} event
     */
    sendGameEvent: function (event) {
        console.log('GAME_EVENT:', event);
        const frame = document.getElementById('GameFrame');
        frame.contentWindow.postMessage(
            { type: "game-cmd", payload: event },
            "*"
        );
    },

    // Method to handle layout logic
    autoLayout: function () {
        const width = $(window).width();
        const height = $(window).height();
        const aspectRatio = width / height;

        // Set app properties
        this.Width = width;
        this.Height = height;

        // Remove all layout classes from the layout container
        $('#Layout').removeClass('FullHorizontal SmallHorizontal FullVertical SmallVertical');

        let layoutClass;

        if (width > height) {
            if (aspectRatio < 1.5) {
                layoutClass = 'SmallHorizontal';
            } else {
                layoutClass = 'FullHorizontal';
            }
        } else {
            const verticalAspectRatio = height / width;
            if (verticalAspectRatio < 1.5) {
                layoutClass = 'SmallVertical';
            } else {
                layoutClass = 'FullVertical';
            }
        }

        // Add the determined layout class
        $('#Layout').addClass(layoutClass);

        // Set the CurrentLayout property
        this.CurrentLayout = layoutClass;

        // Redistribute buttons when layout changes (with a small delay to ensure DOM is updated)
        setTimeout(() => {
            this.ensureButtons();
            this.layoutButtons();
            this.ensureFooterButtons();
            this.layoutFooter();
        }, 10);
    },

    /**
     * Set the currently active brush button by name.
     * Only one brush button is active at any time.
     * @param {string} name
     */
    setBrush: function (name) {
        this.CurrentBrushName = name;
        (this.BrushButtons || []).forEach((entry) => {
            const isActive = entry.name === name;
            entry.button.setActive(isActive);
        });

        // Notify game about brush selection
        this.sendGameEvent('brush_select:' + name.toLowerCase());
    },

    /**
     * Set the active brush coming from the Go game (keyboard, etc.).
     * This only updates UI state, it does NOT send a message back to the game
     * to avoid event echo loops.
     */
    setBrushFromGame: function (name) {
        this.CurrentBrushName = name;
        (this.BrushButtons || []).forEach((entry) => {
            const isActive = entry.name === name;
            entry.button.setActive(isActive);
        });
    },

    /**
     * Ensure the Start/Stop buttons visually reflect the current run state.
     * Exactly one of them is always active.
     */
    ensureRunButtonVisuals: function () {
        const startEntry = (this.Buttons || []).find((entry) => entry.name === 'Start');
        const stopEntry = (this.Buttons || []).find((entry) => entry.name === 'Stop');
        if (!startEntry || !stopEntry) {
            return;
        }

        if (this.IsRunning) {
            startEntry.button.setActive(true);
            stopEntry.button.setActive(false);
        } else {
            startEntry.button.setActive(false);
            stopEntry.button.setActive(true);
        }
    },

    /**
     * Create SpriteButton instances for the 16 main tools if not already created.
     */
    ensureButtons: function () {
        if (this.Buttons && this.Buttons.length) {
            return;
        }

        const buttonNames = [
            'Empty', 'Stone', 'Wood', 'Sand',
            'Water', 'Seed', 'Ant', 'Acid',
            'Fire', 'Ice',
            // right-hand tools in this exact order:
            // gen, erase, size plus, size minus, start, stop
            'Gen', 'Erase', 'SizePlus', 'SizeMinus', 'Start', 'Stop'
        ];
        const brushNames = [
            'Empty', 'Stone', 'Wood', 'Sand',
            'Water', 'Seed', 'Ant', 'Acid',
            'Fire', 'Ice'
        ];

        this.BrushButtons = [];

        this.Buttons = buttonNames.map((name) => {
            const data = spriteData[name];
            if (!data) {
                return null;
            }

            const isBrush = brushNames.includes(name);

            let onClick;

            if (isBrush) {
                onClick = () => {
                    this.setBrush(name);
                };
            } else if (name === 'Start') {
                // Start simulation: mark Start active and Stop inactive.
                onClick = (event, btn) => {
                    if (this.IsRunning) {
                        return;
                    }

                    this.IsRunning = true;
                    this.sendGameEvent('world:start');

                    // Toggle button visuals (exactly one active)
                    this.ensureRunButtonVisuals();
                };
            } else if (name === 'Stop') {
                // Stop simulation: mark Stop active and Start inactive.
                onClick = (event, btn) => {
                    if (!this.IsRunning) {
                        return;
                    }

                    this.IsRunning = false;
                    this.sendGameEvent('world:stop');

                    // Toggle button visuals (exactly one active)
                    this.ensureRunButtonVisuals();
                };
            } else if (name === 'Erase') {
                onClick = (event, btn) => {
                    this.showEraseConfirm();
                    btn.setActive(true);
                };
            } else if (name === 'Gen') {
                onClick = (event, btn) => {
                    this.showGenConfirm();
                    btn.setActive(true);
                };
            } else if (name === 'SizePlus') {
                onClick = (event, btn) => {
                    this.sendGameEvent('brush_size:plus');
                    btn.setActive(true);
                    setTimeout(() => btn.setActive(false), 150);
                };
            } else if (name === 'SizeMinus') {
                onClick = (event, btn) => {
                    this.sendGameEvent('brush_size:minus');
                    btn.setActive(true);
                    setTimeout(() => btn.setActive(false), 150);
                };
            } else {
                onClick = () => {
                    // Placeholder: hook into other tools later
                    console.log(name + ' button clicked');
                };
            }

            const button = new SpriteButton({
                spriteSheetUrl: 'sprite_sheet.png',
                inactiveSprite: data.inactive_sprite,
                activeSprite: data.active_sprite,
                sheetWidth: SPRITE_SHEET_WIDTH,
                sheetHeight: SPRITE_SHEET_HEIGHT,
                onClick
            });

            const entry = { name, button };
            if (isBrush) {
                this.BrushButtons.push(entry);
            }

            return entry;
        }).filter(Boolean);

        // Set default brush selection
        this.setBrush(this.CurrentBrushName || 'Sand');

        // Default world state: not running, Stop button active, Start inactive.
        this.IsRunning = false;
        this.ensureRunButtonVisuals();
    },

    /**
     * Create footer SpriteButtons (Title on the left, GitHub on the right).
     */
    ensureFooterButtons: function () {
        if (this.FooterButtons) {
            return;
        }

        const footerButtons = {};

        const titleData = spriteData.Title;
        if (titleData) {
            footerButtons.title = new SpriteButton({
                spriteSheetUrl: 'sprite_sheet.png',
                inactiveSprite: titleData.inactive_sprite,
                activeSprite: titleData.active_sprite,
                sheetWidth: SPRITE_SHEET_WIDTH,
                sheetHeight: SPRITE_SHEET_HEIGHT
            });
            footerButtons.title.element.classList.add('title-button');
        }

        const githubData = spriteData.GitHub;
        if (githubData) {
            footerButtons.github = new SpriteButton({
                spriteSheetUrl: 'sprite_sheet.png',
                inactiveSprite: githubData.inactive_sprite,
                activeSprite: githubData.active_sprite,
                sheetWidth: SPRITE_SHEET_WIDTH,
                sheetHeight: SPRITE_SHEET_HEIGHT,
                onClick: () => {
                    console.log('GitHub button clicked');
                }
            });
        }

        this.FooterButtons = footerButtons;
    },

    /**
     * Position buttons inside a container as a rows x cols grid,
     * scaling each button to fill its cell while keeping aspect ratio.
     */
    _layoutGrid: function (container, buttons, rows, cols) {
        if (!container || !buttons || !buttons.length) {
            return;
        }

        const rect = container.getBoundingClientRect();
        if (rect.width === 0 || rect.height === 0) {
            return;
        }

        // Account for CSS grid gaps so buttons never overflow the page
        const style = window.getComputedStyle(container);
        const gapX = parseFloat(style.columnGap || style.gap || '0') || 0;
        const gapY = parseFloat(style.rowGap || style.gap || '0') || 0;

        const totalGapX = gapX * (cols - 1);
        const totalGapY = gapY * (rows - 1);

        const cellWidth = (rect.width - totalGapX) / cols;
        const cellHeight = (rect.height - totalGapY) / rows;

        buttons.forEach((entry) => {
            const btn = entry.button || entry;
            const frame = btn.inactiveSprite || btn.activeSprite;
            if (!frame) {
                return;
            }

            const scale = Math.min(cellWidth / frame.w, cellHeight / frame.h);
            btn.setScale(scale);
            container.appendChild(btn.element);
        });
    },

    /**
     * Arrange all tool buttons according to the current layout mode.
     */
    layoutButtons: function () {
        const panelA = document.getElementById('PanelA');
        const panelB = document.getElementById('PanelB');

        [panelA, panelB].forEach((panel) => {
            if (panel) {
                panel.innerHTML = '';
            }
        });

        const allButtons = (this.Buttons || []).map((entry) => entry.button);

        switch (this.CurrentLayout) {
            case 'FullVertical':
                // Single 4x4 grid for all 16 buttons
                this._layoutGrid(panelA, allButtons.slice(0, 16), 4, 4);
                break;
            case 'SmallVertical':
                // Single 2x8 grid at the bottom
                this._layoutGrid(panelA, allButtons.slice(0, 16), 2, 8);
                break;
            case 'FullHorizontal': {
                // Left: first 10 buttons (5x2), Right: remaining 6 (5x2, last 4 cells are empty)
                const leftButtons = allButtons.slice(0, 10);
                const rightButtons = allButtons.slice(10, 16);
                this._layoutGrid(panelA, leftButtons, 5, 2);
                this._layoutGrid(panelB, rightButtons, 5, 2);
                break;
            }
            case 'SmallHorizontal':
                // Single 8x2 grid on the left
                this._layoutGrid(panelA, allButtons.slice(0, 16), 8, 2);
                break;
        }
    },

    /**
     * Footer layout: Title on the left and GitHub button on the right
     */
    layoutFooter: function () {
        const footer = document.getElementById('Footer');
        const panelA = document.getElementById('PanelA');
        const footerTitle = document.getElementById('FooterTitle');
        const footerGitHub = document.getElementById('FooterGitHub');

        if (!footer || !this.FooterButtons) {
            return;
        }

        if (footerTitle) {
            footerTitle.innerHTML = '';
        }
        if (footerGitHub) {
            footerGitHub.innerHTML = '';
        }

        const rect = footer.getBoundingClientRect();
        if (rect.height === 0) {
            return;
        }

        const availableHeight = rect.height - 4; // small padding

        // In FullVertical mode, keep the title button's width aligned
        // with the button panel (PanelA) width so it never "sticks out"
        // under the game area on very small screens.
        let maxTitleWidth = null;
        if (this.CurrentLayout === 'FullVertical' && panelA) {
            const panelRect = panelA.getBoundingClientRect();
            if (panelRect.width > 0) {
                maxTitleWidth = panelRect.width;
            }
        }

        const { title, github } = this.FooterButtons;

        if (title && footerTitle) {
            const frame = title.inactiveSprite || title.activeSprite;
            if (frame) {
                const heightScale = availableHeight / frame.h;
                let scale = heightScale;

                if (maxTitleWidth != null) {
                    const widthScale = maxTitleWidth / frame.w;
                    scale = Math.min(heightScale, widthScale);
                }

                title.setScale(scale);
            }
            title.mount(footerTitle);
        }

        if (github && footerGitHub) {
            const frame = github.inactiveSprite || github.activeSprite;
            if (frame) {
                const scale = availableHeight / frame.h;
                github.setScale(scale);
            }
            github.mount(footerGitHub);
        }
    },

    /**
     * Internal helper to show a confirmation dialog inside the GameWindow.
     * Any subsequent click on the page counts as an answer:
     * - click on the primary button => actual action (erase/gen)
     * - click anywhere else => treated as 'no'
     * @param {'erase'|'gen'} type
     */
    showConfirmDialog: function (type) {
        const overlay = document.getElementById('ConfirmOverlay');
        const dialog = document.getElementById('ConfirmDialog');
        const textEl = document.getElementById('ConfirmText');
        const primaryBtn = document.getElementById('ConfirmPrimary');
        const cancelBtn = document.getElementById('ConfirmCancel');
        const gameWindow = document.getElementById('GameWindow');

        if (!overlay || !dialog || !textEl || !primaryBtn || !cancelBtn || !gameWindow) {
            return;
        }

        // If a modal is already active, ignore additional requests
        if (this.ActiveModalType) {
            return;
        }
        this.ActiveModalType = type;

        const rect = gameWindow.getBoundingClientRect();
        const dialogWidth = Math.max(0, rect.width - 10);
        const dialogHeight = (2 / 3) * rect.height;

        dialog.style.width = `${dialogWidth}px`;
        dialog.style.height = `${dialogHeight}px`;

        // Scale font relative to dialog size
        const baseSize = Math.min(dialogHeight / 6, dialogWidth / 18);
        const fontSize = Math.max(8, Math.min(18, baseSize));
        dialog.style.fontSize = `${fontSize}px`;

        // Configure text and primary button appearance
        primaryBtn.classList.remove('confirm-primary-erase', 'confirm-primary-gen');

        if (type === 'erase') {
            textEl.textContent = 'do you really want to erase the current state, and start a new empty world?';
            primaryBtn.textContent = 'erase';
            primaryBtn.classList.add('confirm-primary-erase');
        } else if (type === 'gen') {
            textEl.textContent = 'do you really want to generate a new world?';
            primaryBtn.textContent = 'gen';
            primaryBtn.classList.add('confirm-primary-gen');
        }

        cancelBtn.textContent = 'no';

        // Highlight the corresponding tool button while the dialog is open
        const toolName = type === 'erase' ? 'Erase' : 'Gen';
        const toolEntry = (this.Buttons || []).find((entry) => entry.name === toolName);
        if (toolEntry) {
            toolEntry.button.setActive(true);
        }

        overlay.classList.add('visible');

        // Any subsequent click on the document counts as a response.
        // Attach after a tick so we don't catch the opening click.
        setTimeout(() => {
            const handler = (event) => {
                const primary = document.getElementById('ConfirmPrimary');
                const isPrimary = primary && primary.contains(event.target);
                const activeType = this.ActiveModalType;

                if (isPrimary && activeType === 'erase') {
                    this.sendGameEvent('world:erase');
                } else if (isPrimary && activeType === 'gen') {
                    this.sendGameEvent('world:gen');
                }

                this.hideConfirm();

                event.stopPropagation();
                event.preventDefault();
            };

            document.addEventListener('click', handler, { capture: true, once: true });
        }, 0);
    },

    /**
     * Show erase confirmation popup.
     */
    showEraseConfirm: function () {
        this.showConfirmDialog('erase');
    },

    /**
     * Show generate confirmation popup.
     */
    showGenConfirm: function () {
        this.showConfirmDialog('gen');
    },

    /**
     * Hide any confirmation popup.
     */
    hideConfirm: function () {
        const overlay = document.getElementById('ConfirmOverlay');
        if (overlay) {
            overlay.classList.remove('visible');
        }

        // Reset tool button highlight
        if (this.ActiveModalType === 'erase' || this.ActiveModalType === 'gen') {
            const toolName = this.ActiveModalType === 'erase' ? 'Erase' : 'Gen';
            const toolEntry = (this.Buttons || []).find((entry) => entry.name === toolName);
            if (toolEntry) {
                toolEntry.button.setActive(false);
            }
        }

        this.ActiveModalType = null;
    },

    // Initialize method to set up event listeners
    initialize: async function () {

        this.autoLayout();

        // Set up event listeners
        $(window).on('resize', () => {
            this.autoLayout();
        });

        $(window).on('orientationchange', () => {
            // Add a small delay to ensure dimensions are updated
            setTimeout(() => {
                this.autoLayout();
            }, 100);
        });

        // Listen for messages coming from the WASM game iframe.
        // Expected shape: { type: 'ui-cmd', payload: '...' }
        window.addEventListener('message', (event) => {
            const data = event.data;
            if (!data || data.type !== 'ui-cmd' || typeof data.payload !== 'string') {
                return;
            }

            const payload = data.payload;

            // Brush selection coming from Go: \"brush_select:<name>\"
            if (payload.startsWith('brush_select:')) {
                const parts = payload.split(':');
                const key = (parts[1] || '').trim();
                if (!key) {
                    return;
                }

                // Convert \"sand\" -> \"Sand\" etc. to match our button naming.
                const brushName = key.charAt(0).toUpperCase() + key.slice(1);
                this.setBrushFromGame(brushName);
            }
        });
    }
};

// Initialize the app when the DOM is ready (modern jQuery shorthand)
$(function () {
    App.initialize();
});
