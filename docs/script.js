const SMALL_FULL_THRESHOLD = 1.5;
const SPRITE_SHEET_WIDTH = 190;
const SPRITE_SHEET_HEIGHT = 190;
const SPRITE_SIZE = 26;

const SPRITE_SHEET_DATA = {
    Empty: {
        Inactive: {
            x: 1,
            y: 1,
        },
        Active: {
            x: 28,
            y: 1,
        }
    },
    Stone: {
        Inactive: {
            x: 55,
            y: 1,
        },
        Active: {
            x: 82,
            y: 1,
        }
    },
    Sand: {
        Inactive: {
            x: 109,
            y: 1,
        },
        Active: {
            x: 136,
            y: 1,
        }
    },
    Water: {
        Inactive: {
            x: 163,
            y: 1,
        },
        Active: {
            x: 1,
            y: 28,
        }
    },
    Seed: {
        Inactive: {
            x: 28,
            y: 28,
        },
        Active: {
            x: 55,
            y: 28,
        }
    },
    Ant: {
        Inactive: {
            x: 82,
            y: 28,
        },
        Active: {
            x: 109,
            y: 28,
        }
    },
    Wasp: {
        Inactive: {
            x: 136,
            y: 28,
        },
        Active: {
            x: 163,
            y: 28,
        }
    },
    Acid: {
        Inactive: {
            x: 1,
            y: 55,
        },
        Active: {
            x: 28,
            y: 55,
        }
    },
    Fire: {
        Inactive: {
            x: 55,
            y: 55,
        },
        Active: {
            x: 82,
            y: 55,
        }
    },
    Ice: {
        Inactive: {
            x: 109,
            y: 55,
        },
        Active: {
            x: 136,
            y: 55,
        }
    },
    Start: {
        Inactive: {
            x: 163,
            y: 55,
        },
        Active: {
            x: 1,
            y: 82,
        }
    },
    Stop: {
        Inactive: {
            x: 28,
            y: 82,
        },
        Active: {
            x: 55,
            y: 82,
        }
    },
    Erase: {
        Inactive: {
            x: 82,
            y: 82,
        },
        Active: {
            x: 109,
            y: 82,
        }
    },
    Gen: {
        Inactive: {
            x: 136,
            y: 82,
        },
        Active: {
            x: 163,
            y: 82,
        }
    },
    Size1: {
        Inactive: {
            x: 1,
            y: 109,
        },
        Active: {
            x: 28,
            y: 109,
        }
    },
    Size2: {
        Inactive: {
            x: 55,
            y: 109,
        },
        Active: {
            x: 82,
            y: 109,
        }
    },
    Size3: {
        Inactive: {
            x: 109,
            y: 109,
        },
        Active: {
            x: 136,
            y: 109,
        }
    },
    Size4: {
        Inactive: {
            x: 163,
            y: 109,
        },
        Active: {
            x: 1,
            y: 136,
        }
    },
    Size5: {
        Inactive: {
            x: 28,
            y: 136,
        },
        Active: {
            x: 55,
            y: 136,
        }
    },
    Menu: {
        Inactive: {
            x: 82,
            y: 136,
        },
        Active: {
            x: 109,
            y: 136,
        }
    },
    Info: {
        Inactive: {
            x: 109,
            y: 136,
        },
        Active: {
            x: 136,
            y: 136,
        }
    },
    Options: {
        Inactive: {
            x: 163,
            y: 136,
        },
        Active: {
            x: 1,
            y: 163,
        }
    }
};

const DIALOG_DATA = {
    Erase: {
        Color: "red",
        Message: "do you really want to erase the world?",
        YesAnswer: "erase",
        NoAnswer: "cancel",
    },
    Gen: {
        Color: "green",
        Message: "do you really want to generate a new world?",
        YesAnswer: "gen",
        NoAnswer: "cancel",
    }
}

const MENU_ITEMS = [
    {
        Text: "full screen: off",
        Action: () => { switchFullscreen(); }
    },
    {
        Text: "debug mode: off",
        Action: () => {
            DebugMode = !DebugMode;
            MENU_ITEMS[1].Text = "debug mode: " + (DebugMode ? "on" : "off");
            sendSiteEventToGame("world:debug:" + (DebugMode ? "on" : "off"));
        }
    },
    {
        Text: "auto rotate world: on",
        Action: () => {
            AutoRotateWorld = !AutoRotateWorld;
            MENU_ITEMS[2].Text = "auto rotate world: " + (AutoRotateWorld ? "on" : "off");
        }
    },
    {
        Text: "rotate the world CW",
        Action: () => {
            sendSiteEventToGame("world:rotate_cw");
        }
    },
    {
        Text: "rotate the world CCW",
        Action: () => {
            sendSiteEventToGame("world:rotate_ccw");
        }
    },
    {
        Text: "flip the world",
        Action: () => {
            sendSiteEventToGame("world:rotate_cw");
            sendSiteEventToGame("world:rotate_cw");
        }
    },
];

const Buttons = {
    Empty: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Empty],
        ButtonAction: () => { setActiveBrush("Empty"); },
    },
    Stone: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Stone],
        ButtonAction: () => { setActiveBrush("Stone"); },
    },
    Sand: {
        ButtonType: "Brush",
        IsActive: true,
        Sprites: [SPRITE_SHEET_DATA.Sand],
        ButtonAction: () => { setActiveBrush("Sand"); },
    },
    Water: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Water],
        ButtonAction: () => { setActiveBrush("Water"); },
    },
    Seed: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Seed],
        ButtonAction: () => { setActiveBrush("Seed"); },
    },
    Ant: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Ant],
        ButtonAction: () => { setActiveBrush("Ant"); },
    },
    Wasp: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Wasp],
        ButtonAction: () => { setActiveBrush("Wasp"); },
    },
    Acid: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Acid],
        ButtonAction: () => { setActiveBrush("Acid"); },
    },
    Fire: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Fire],
        ButtonAction: () => { setActiveBrush("Fire"); },
    },
    Ice: {
        ButtonType: "Brush",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Ice],
        ButtonAction: () => { setActiveBrush("Ice"); },
    },
    Play: {
        ButtonType: "World",
        Sprites: [
            SPRITE_SHEET_DATA.Start,
            SPRITE_SHEET_DATA.Stop,
        ],
        ButtonAction: () => {
            IsRunning = !IsRunning;
            setButtonActivity("Play", true);

            sendSiteEventToGame("world:" + (IsRunning ? "start" : "stop"));

            setTimeout(() => {
                setButtonActivity("Play", false);
            }, 200);
        },
    },
    Size: {
        ButtonType: "World",
        Sprites: [
            SPRITE_SHEET_DATA.Size1,
            SPRITE_SHEET_DATA.Size2,
            SPRITE_SHEET_DATA.Size3,
            SPRITE_SHEET_DATA.Size4,
            SPRITE_SHEET_DATA.Size5,
        ],
        ButtonAction: () => {
            const brush_scale = {
                1: 8,
                2: 14,
                3: 20,
                4: 26,
                5: 32
            }

            setButtonActivity("Size", true);

            BrushSize += 1;
            if (BrushSize > 5) BrushSize = 1;

            sendSiteEventToGame("brush_size:" + brush_scale[BrushSize]);

            setTimeout(() => {
                setButtonActivity("Size", false);
            }, 200);
        },
    },
    Erase: {
        ButtonType: "World",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Erase],
        ButtonAction: () => {
            openConfirmDialog("Erase");
        },
    },
    Gen: {
        ButtonType: "World",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Gen],
        ButtonAction: () => {
            openConfirmDialog("Gen");
        },
    },
    Menu: {
        ButtonType: "World",
        IsActive: false,
        Sprites: [SPRITE_SHEET_DATA.Menu],
        ButtonAction: () => {
            openMenuPopup();
        },
    },
};

const Layouts = {
    FullVertical: {
        PanelA: {
            Columns: 4,
            Rows: 4,
            Buttons: [
                "Size",
                "Play",
                "Gen",
                "Erase",
                "Empty",
                "Stone",
                "Sand",
                "Water",
                "Seed",
                "Ant",
                "Wasp",
                "Acid",
                "Fire",
                "Ice",
                "NONE",
                "Menu",
            ],
        },
        PanelB: null,
        MainContainer: {
            FlexDirection: "column",
            Order: ["GameWindow", "ButtonPanelA"],
        },
        GameWindow: {
            Width: "100cqw",
            Height: "100cqw",
            Flex: "0 0 auto",
        },
        GameFrame: {
            Border: "border-full-vertical",
        }
    },
    SmallVertical: {
        PanelA: {
            Columns: 8,
            Rows: 2,
            Buttons: [
                "Empty",
                "Stone",
                "Sand",
                "Water",
                "Seed",
                "Gen",
                "Erase",
                "Size",
                "Ant",
                "Wasp",
                "Acid",
                "Fire",
                "Ice",
                "Play",
                "NONE",
                "Menu",
            ],
        },
        PanelB: null,
        MainContainer: {
            FlexDirection: "column",
            Order: ["GameWindow", "ButtonPanelA"],
        },
        GameWindow: {
            Flex: "1 1 auto",
        },
        GameFrame: {
            Border: "border-small-vertical",
        }
    },
    FullHorizontal: {
        PanelA: {
            Columns: 2,
            Rows: 5,
            Buttons: [
                "Empty",
                "Stone",
                "Sand",
                "Water",
                "Seed",
                "Ant",
                "Wasp",
                "Acid",
                "Fire",
                "Ice",
            ],
        },
        PanelB: {
            Columns: 2,
            Rows: 5,
            Buttons: [
                "Play",
                "Size",
                "Gen",
                "Erase",
                "NONE",
                "NONE",
                "NONE",
                "NONE",
                "NONE",
                "Menu",
            ],
        },
        MainContainer: {
            FlexDirection: "row",
            Order: ["ButtonPanelA", "GameWindow", "ButtonPanelB"],
        },
        GameWindow: {
            Width: "100cqh",
            Height: "100cqh",
            Flex: "0 0 auto",
        },
        GameFrame: {
            Border: "border-full-horizontal",
        }
    },
    SmallHorizontal: {
        PanelA: {
            Columns: 2,
            Rows: 8,
            Buttons: [
                "Empty",
                "Stone",
                "Sand",
                "Water",
                "Seed",
                "Ant",
                "Wasp",
                "Acid",
                "Fire",
                "Ice",
                "Size",
                "Play",
                "Gen",
                "Erase",
                "Menu",
                "NONE",
            ],
        },
        PanelB: null,
        MainContainer: {
            FlexDirection: "row",
            Order: ["ButtonPanelA", "GameWindow"],
        },
        GameWindow: {
            Flex: "1 1 auto",
        },
        GameFrame: {
            Border: "border-small-horizontal",
        }
    },
};

// Reference to HTML elements
const MainContainer = document.getElementById("MainContainer");
const GameWindow = document.getElementById("GameWindow");
const GameContainer = document.getElementById("GameContainer");
const GameFrame = document.getElementById("GameFrame");
const ButtonPanelA = document.getElementById("ButtonPanelA");
const ButtonPanelB = document.getElementById("ButtonPanelB");
const PopUp = document.getElementById("PopUp");

// Game state
let CurrentLayout = "NoLayout";
let IsRunning = true;
let DebugMode = false;
let FullScreenMode = false;
let AutoRotateWorld = true;
let ActiveBrush = "Sand";
let BrushSize = 2;
let ActivePopUp = null;
let ActivePopUpWindowEl = null;
let RefreshMenu = null;
let LastAngle = screen.orientation?.angle ?? 0;

// pass these string events to the game inside the WebAssembly game frame
function sendSiteEventToGame(event) {
    console.log("SITE_EVENT:", event);
    GameFrame.contentWindow.postMessage(
        { type: "site-event", payload: event },
        "*"
    );
}

// handle string events coming from the game (triggered by keyboard and mouse)
function handleGameEvent(event) {
    const data = event.data;
    if (!data || data.type !== "game-event" || typeof data.payload !== "string") {
        return;
    }

    const payload = data.payload;

    console.log("GAME_EVENT:", payload);

    if (payload.startsWith("game:fullscreen")) {
        switchFullscreen();
    }

    if (payload.startsWith("world:start")) {
        IsRunning = true;
        renderButtons();
    }

    if (payload.startsWith("world:stop")) {
        IsRunning = false;
        renderButtons();
    }

    if (payload.startsWith("world:debug:")) {
        const parts = payload.split(":");
        const key = (parts[2] || "").trim();
        DebugMode = (key === "on");
        MENU_ITEMS[1].Text = "debug mode: " + (DebugMode ? "on" : "off");
    }

    if (payload.startsWith("brush_select:")) {
        const parts = payload.split(":");
        const key = (parts[1] || "").trim();
        setActiveBrush(key);
    }

    if (payload.startsWith("brush_size:")) {
        const parts = payload.split(":");
        const key = (parts[1] || "").trim();
        const num = parseInt(key);

        const oldBrushSize = BrushSize;

        if (num < 9) {
            BrushSize = 1;
        } else if (num < 15) {
            BrushSize = 2;
        } else if (num < 21) {
            BrushSize = 3;
        } else if (num < 27) {
            BrushSize = 4;
        } else {
            BrushSize = 5;
        }

        if (oldBrushSize !== BrushSize) {
            renderButtons();
        }
    }
}

function switchFullscreen() {
    if (document.fullscreenElement) {
        document.exitFullscreen();
        FullScreenMode = false;
    } else {
        MainContainer.requestFullscreen();
        FullScreenMode = true;
    }
    MENU_ITEMS[0].Text = "full screen: " + (FullScreenMode ? "on" : "off");
}

// Detect the current layout based on aspect ratio
function detectLayoutMode() {
    const rect = MainContainer.getBoundingClientRect();
    const W = rect.width;
    const H = rect.height;

    if (W <= 0 || H <= 0) return "FullHorizontal";

    const isHorizontal = (W > H);

    const ratio = isHorizontal ? (W / H) : (H / W);

    if (ratio < SMALL_FULL_THRESHOLD) {
        return isHorizontal ? "SmallHorizontal" : "SmallVertical";
    }

    return isHorizontal ? "FullHorizontal" : "FullVertical";
}


function setSpriteBackgroundStretch(el, sx, sy) {
    const { width: elW, height: elH } = el.getBoundingClientRect();
    const scaleX = elW / SPRITE_SIZE;
    const scaleY = elH / SPRITE_SIZE;

    // Scale the whole sheet so the sprite rect maps to the element size
    el.style.backgroundSize = `${SPRITE_SHEET_WIDTH * scaleX}px ${SPRITE_SHEET_HEIGHT * scaleY}px`;

    // Shift sheet so the sprite's top-left aligns with the element's top-left
    el.style.backgroundPosition = `${-sx * scaleX}px ${-sy * scaleY}px`;
}

// build the buttons grid according to the current layout
function buildGrid(el, layout) {
    el.textContent = "";
    el.style.gridTemplateColumns = `repeat(${layout.Columns}, 1fr)`;
    el.style.gridTemplateRows = `repeat(${layout.Rows}, 1fr)`;

    for (let i = 0; i < layout.Columns * layout.Rows; i++) {
        const btnContainer = document.createElement("div");
        btnContainer.classList.add("btn-container");

        const btnName = layout.Buttons[i];
        if (btnName !== "NONE") {
            const btn = document.createElement("div");
            btn.classList.add("btn");
            btn.id = btnName + "Button";
            btn.addEventListener("click", Buttons[btnName].ButtonAction);
            btnContainer.appendChild(btn);
        }

        el.appendChild(btnContainer);
    }
}

function setCssOrClear(el, prop, value) {
    el.style[prop] = (value === undefined || value === null) ? "" : value;
}

function applyPanelConstraints(layout) {
    // Clear previous constraints (important when switching layouts)
    ButtonPanelA.style.flex = "";
    ButtonPanelA.style.width = "";
    ButtonPanelA.style.height = "";
    ButtonPanelA.style.minWidth = "";
    ButtonPanelA.style.minHeight = "";

    // Also clear any clamp overrides; we'll re-apply if needed
    // (Do not clear for Full layouts; layout-change block sets those)
    if (layout === "SmallVertical" || layout === "SmallHorizontal") {
        GameWindow.style.width = "";
        GameWindow.style.height = "";
        GameWindow.style.flex = "1 1 auto";
    }

    const minStrip = 4 * SPRITE_SIZE;

    const mc = MainContainer.getBoundingClientRect();
    const W = mc.width;
    const H = mc.height;
    if (W <= 0 || H <= 0) return;

    if (layout === "SmallVertical") {
        // Default: panel is minimum; game takes the rest
        // If game would be taller than wide, clamp game height to width
        // and give leftover height to the panel.
        const clamp = (W < (H - minStrip));

        const gameH = clamp ? Math.floor(W) : null;            // clamp to square
        const stripH = clamp ? Math.floor(H - W) : minStrip;   // absorb leftover

        ButtonPanelA.style.flex = `0 0 ${stripH}px`;
        ButtonPanelA.style.height = `${stripH}px`;
        ButtonPanelA.style.minHeight = `${stripH}px`;

        GameWindow.style.width = "100%";
        if (clamp) {
            GameWindow.style.flex = `0 0 ${gameH}px`;
            GameWindow.style.height = `${gameH}px`;
        } else {
            GameWindow.style.flex = "1 1 auto";
            GameWindow.style.height = "";
        }
    }

    if (layout === "SmallHorizontal") {
        // Default: panel is minimum; game takes the rest
        // If game would be wider than tall, clamp game width to height
        // and give leftover width to the panel.
        const clamp = (H < (W - minStrip));

        const gameW = clamp ? Math.floor(H) : null;            // clamp to square
        const stripW = clamp ? Math.floor(W - H) : minStrip;   // absorb leftover

        ButtonPanelA.style.flex = `0 0 ${stripW}px`;
        ButtonPanelA.style.width = `${stripW}px`;
        ButtonPanelA.style.minWidth = `${stripW}px`;

        GameWindow.style.height = "100%";
        if (clamp) {
            GameWindow.style.flex = `0 0 ${gameW}px`;
            GameWindow.style.width = `${gameW}px`;
        } else {
            GameWindow.style.flex = "1 1 auto";
            GameWindow.style.width = "";
        }
    }
}

function applyGameSquare() {
    const gw = GameWindow.getBoundingClientRect();
    const gwW = gw.width;
    const gwH = gw.height;

    if (gwW <= 0 || gwH <= 0) return;

    const size = Math.max(0, Math.floor(Math.min(gwW, gwH)));

    GameContainer.style.width = `${size}px`;
    GameContainer.style.height = `${size}px`;

    GameFrame.style.width = "100%";
    GameFrame.style.height = "100%";
}

function setActiveBrush(brushName) {
    if (brushName === ActiveBrush) return;
    console.log(`ActiveBrush is set from ${ActiveBrush} to ${brushName}`);
    ActiveBrush = brushName;
    Object.keys(Buttons).forEach(btnName => {
        const btn = Buttons[btnName];
        if (btn.ButtonType === "Brush") {
            btn.IsActive = (btnName === brushName);
        }
    });
    renderButtons();
    sendSiteEventToGame("brush_select:" + ActiveBrush.toLowerCase());
}

function renderButtons() {
    Object.keys(Buttons).forEach(btnName => {
        const btn = Buttons[btnName];
        const el = document.getElementById(btnName + "Button");
        if (el) {
            let sprIdx = 0;
            if (btnName == "Size") {
                sprIdx = BrushSize - 1;
            }
            if (btnName == "Play") {
                if (IsRunning) {
                    sprIdx = 1;
                }
            }
            if (btn.IsActive) {
                setSpriteBackgroundStretch(el, btn.Sprites[sprIdx].Active.x, btn.Sprites[sprIdx].Active.y);
            } else {
                setSpriteBackgroundStretch(el, btn.Sprites[sprIdx].Inactive.x, btn.Sprites[sprIdx].Inactive.y);
            }
        }
    });
}

function setButtonActivity(btnName, isActive) {
    Buttons[btnName].IsActive = isActive;
    renderButtons();
}

function closePopUpAndDeactivate(openerBtnName) {
    hidePopUp();
    setButtonActivity(openerBtnName, false);
}

function positionPopUpWindow(win) {
    const mc = MainContainer.getBoundingClientRect();
    const gf = GameFrame.getBoundingClientRect();

    const w = Math.floor(gf.width * 0.9);
    const h = Math.floor(gf.height * 0.9);

    const left = Math.floor((gf.left - mc.left) + (gf.width - w) / 2);
    const top = Math.floor((gf.top - mc.top) + (gf.height - h) / 2);

    win.style.width = `${w}px`;
    win.style.height = `${h}px`;
    win.style.left = `${left}px`;
    win.style.top = `${top}px`;
}

function hidePopUp() {
    PopUp.classList.add("hidden");
    PopUp.innerHTML = "";
    PopUp.onclick = null;
    ActivePopUp = null;

    ActivePopUpWindowEl = null;
    RefreshMenu = null;
}

function showPopUpShell(borderColor, onOutsideClick) {
    PopUp.classList.remove("hidden");
    PopUp.innerHTML = "";

    // Clicking anywhere NOT on the window counts as outside (because overlay covers MainContainer)
    PopUp.onclick = (e) => {
        if (e.target === PopUp) {
            onOutsideClick?.();
        }
    };

    const win = document.createElement("div");
    win.className = "popup-window";
    win.style.borderColor = borderColor;

    // Prevent inside clicks from becoming outside clicks
    win.addEventListener("click", (e) => e.stopPropagation());

    PopUp.appendChild(win);

    ActivePopUpWindowEl = win;
    positionPopUpWindow(win);

    return win;
}


function openConfirmDialog(openerBtnName) {
    if (ActivePopUp) return; // keep current active until answered

    const cfg = DIALOG_DATA[openerBtnName];

    ActivePopUp = { Type: openerBtnName };

    setButtonActivity(openerBtnName, true);

    const win = showPopUpShell(cfg.Color, () => {
        // outside click counts as No
        closePopUpAndDeactivate(openerBtnName);
    });

    const content = document.createElement("div");
    content.className = "popup-content";

    const msg = document.createElement("div");
    msg.textContent = cfg.Message; // white by default
    content.appendChild(msg);

    const actions = document.createElement("div");
    actions.className = "popup-actions";

    const yesBtn = document.createElement("button");
    yesBtn.className = "popup-btn";
    yesBtn.textContent = cfg.YesAnswer;
    yesBtn.style.color = cfg.Color;        // YesAnswer text color
    yesBtn.style.borderColor = cfg.Color;  // YesAnswer border color

    yesBtn.addEventListener("click", () => {
        // Send event to the game (adjust message format if you prefer)
        sendSiteEventToGame("world:" + cfg.YesAnswer);
        closePopUpAndDeactivate(openerBtnName);
    });

    const noBtn = document.createElement("button");
    noBtn.className = "popup-btn";
    noBtn.textContent = cfg.NoAnswer;
    noBtn.style.color = "white";
    noBtn.style.borderColor = "white";

    noBtn.addEventListener("click", () => {
        closePopUpAndDeactivate(openerBtnName);
    });

    actions.appendChild(yesBtn);
    actions.appendChild(noBtn);

    win.appendChild(content);
    win.appendChild(actions);
}

function openMenuPopup() {
    if (ActivePopUp) return; // keep current active until closed

    ActivePopUp = { Type: "Menu" };
    setButtonActivity("Menu", true);

    const win = showPopUpShell("white", () => {
        // outside click closes menu
        closePopUpAndDeactivate("Menu");
    });

    const menuWrap = document.createElement("div");
    menuWrap.className = "popup-menu";

    function rebuildMenuButtons() {
        menuWrap.innerHTML = "";
        for (const item of MENU_ITEMS) {
            const b = document.createElement("button");
            b.className = "popup-btn popup-menu-btn";
            b.textContent = item.Text;
            b.addEventListener("click", () => {
                item.Action?.();
                rebuildMenuButtons(); // text may change after action
            });
            menuWrap.appendChild(b);
        }
    }

    RefreshMenu = rebuildMenuButtons;

    rebuildMenuButtons();
    win.appendChild(menuWrap);
}

function normAngle(a) {
    a = ((a % 360) + 360) % 360;

    return a;
}

function rotateWorld() {
    const newAngle = screen.orientation.angle;

    if (AutoRotateWorld) {
        const delta = (normAngle(newAngle) - normAngle(LastAngle) + 360) % 360;

        if (delta === 90) {
            sendSiteEventToGame("world:rotate_ccw");
        }

        if (delta === 180) {
            sendSiteEventToGame("world:rotate_cw");
            sendSiteEventToGame("world:rotate_cw");
        }

        if (delta === 270) {
            sendSiteEventToGame("world:rotate_cw");
        }
    }

    LastAngle = newAngle;
}

function autoLayout() {
    const layout = detectLayoutMode();

    if (CurrentLayout !== layout) {
        console.log(`Layout changed from ${CurrentLayout} to ${layout}`);
        CurrentLayout = layout;

        const newLayout = Layouts[layout];

        // Rebuild grids
        buildGrid(ButtonPanelA, newLayout.PanelA);
        if (newLayout.PanelB) {
            buildGrid(ButtonPanelB, newLayout.PanelB);
            ButtonPanelB.classList.remove("hidden");
        } else {
            ButtonPanelB.classList.add("hidden");
        }

        // Flex direction + order (reset first to avoid stale orders)
        MainContainer.style.flexDirection = newLayout.MainContainer.FlexDirection;

        ["ButtonPanelA", "GameWindow", "ButtonPanelB"].forEach(id => {
            document.getElementById(id).style.order = "";
        });

        for (let i = 0; i < newLayout.MainContainer.Order.length; i++) {
            document.getElementById(newLayout.MainContainer.Order[i]).style.order = i + 1;
        }

        // GameWindow sizing rules (only set if present; otherwise clear)
        setCssOrClear(GameWindow, "width", newLayout.GameWindow.Width);
        setCssOrClear(GameWindow, "height", newLayout.GameWindow.Height);
        setCssOrClear(GameWindow, "flex", newLayout.GameWindow.Flex);

        GameFrame.classList.remove(
            "border-full-vertical",
            "border-small-vertical",
            "border-full-horizontal",
            "border-small-horizontal"
        );

        GameFrame.classList.add(newLayout.GameFrame.Border);
    }

    // These must run even when layout mode didn't change (resizes)
    applyPanelConstraints(CurrentLayout);
    applyGameSquare();
    renderButtons();
    rotateWorld();


    if (!PopUp.classList.contains("hidden") && ActivePopUpWindowEl) {
        positionPopUpWindow(ActivePopUpWindowEl);
    }

}

window.addEventListener("load", async () => {
    // Listen for any change that can influence layout
    window.addEventListener("resize", autoLayout);
    window.addEventListener("orientationchange", autoLayout);
    if (screen.orientation) {
        screen.orientation.addEventListener("change", autoLayout);
    }

    // Listen for game events
    window.addEventListener("message", handleGameEvent);

    // Call autoLayout once for the initial setup
    autoLayout();
});