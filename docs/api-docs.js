const {
  Document,
  Packer,
  Paragraph,
  TextRun,
  Table,
  TableRow,
  TableCell,
  HeadingLevel,
  AlignmentType,
  BorderStyle,
  WidthType,
  ShadingType,
  LevelFormat,
  PageNumber,
  Footer,
  PageBreak,
} = require("docx");
const fs = require("fs");

// ── Helpers ───────────────────────────────────────────────────
const COLORS = {
  primary: "1E3A5F",
  accent: "2E86AB",
  green: "27AE60",
  red: "E74C3C",
  orange: "F39C12",
  purple: "8E44AD",
  lightBg: "F0F4F8",
  headerBg: "1E3A5F",
  codeBg: "F7F9FC",
  border: "D0D7DE",
  text: "24292F",
  muted: "656D76",
};

const METHOD_COLORS = {
  GET: "27AE60",
  POST: "2E86AB",
  PUT: "F39C12",
  PATCH: "8E44AD",
  DELETE: "E74C3C",
};

const border = { style: BorderStyle.SINGLE, size: 1, color: COLORS.border };
const borders = { top: border, bottom: border, left: border, right: border };
const noBorder = { style: BorderStyle.NONE, size: 0, color: "FFFFFF" };
const noBorders = {
  top: noBorder,
  bottom: noBorder,
  left: noBorder,
  right: noBorder,
};

function h1(text) {
  return new Paragraph({
    heading: HeadingLevel.HEADING_1,
    spacing: { before: 360, after: 120 },
    border: {
      bottom: {
        style: BorderStyle.SINGLE,
        size: 4,
        color: COLORS.accent,
        space: 6,
      },
    },
    children: [
      new TextRun({
        text,
        font: "Arial",
        size: 36,
        bold: true,
        color: COLORS.primary,
      }),
    ],
  });
}

function h2(text) {
  return new Paragraph({
    heading: HeadingLevel.HEADING_2,
    spacing: { before: 280, after: 80 },
    children: [
      new TextRun({
        text,
        font: "Arial",
        size: 28,
        bold: true,
        color: COLORS.primary,
      }),
    ],
  });
}

function h3(text) {
  return new Paragraph({
    heading: HeadingLevel.HEADING_3,
    spacing: { before: 200, after: 60 },
    children: [
      new TextRun({
        text,
        font: "Arial",
        size: 24,
        bold: true,
        color: COLORS.accent,
      }),
    ],
  });
}

function p(text, opts = {}) {
  return new Paragraph({
    spacing: { before: 60, after: 60 },
    children: [
      new TextRun({
        text,
        font: "Arial",
        size: 22,
        color: COLORS.text,
        ...opts,
      }),
    ],
  });
}

function muted(text) {
  return new Paragraph({
    spacing: { before: 40, after: 40 },
    children: [
      new TextRun({
        text,
        font: "Arial",
        size: 20,
        color: COLORS.muted,
        italics: true,
      }),
    ],
  });
}

function spacer() {
  return new Paragraph({
    spacing: { before: 80, after: 80 },
    children: [new TextRun("")],
  });
}

function code(text) {
  return new Paragraph({
    spacing: { before: 40, after: 40 },
    indent: { left: 360 },
    children: [
      new TextRun({
        text,
        font: "Courier New",
        size: 18,
        color: COLORS.primary,
      }),
    ],
  });
}

function methodBadge(method, path, description) {
  const color = METHOD_COLORS[method] || COLORS.accent;
  return new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: [1000, 3200, 5160],
    borders: noBorders,
    rows: [
      new TableRow({
        children: [
          new TableCell({
            borders: noBorders,
            width: { size: 1000, type: WidthType.DXA },
            shading: { fill: color, type: ShadingType.CLEAR },
            margins: { top: 60, bottom: 60, left: 120, right: 120 },
            verticalAlign: "center",
            children: [
              new Paragraph({
                alignment: AlignmentType.CENTER,
                children: [
                  new TextRun({
                    text: method,
                    font: "Arial",
                    size: 18,
                    bold: true,
                    color: "FFFFFF",
                  }),
                ],
              }),
            ],
          }),
          new TableCell({
            borders: noBorders,
            width: { size: 3200, type: WidthType.DXA },
            shading: { fill: COLORS.codeBg, type: ShadingType.CLEAR },
            margins: { top: 60, bottom: 60, left: 160, right: 120 },
            children: [
              new Paragraph({
                children: [
                  new TextRun({
                    text: path,
                    font: "Courier New",
                    size: 20,
                    bold: true,
                    color: COLORS.primary,
                  }),
                ],
              }),
            ],
          }),
          new TableCell({
            borders: noBorders,
            width: { size: 5160, type: WidthType.DXA },
            shading: { fill: "FFFFFF", type: ShadingType.CLEAR },
            margins: { top: 60, bottom: 60, left: 200, right: 120 },
            children: [
              new Paragraph({
                children: [
                  new TextRun({
                    text: description,
                    font: "Arial",
                    size: 20,
                    color: COLORS.muted,
                  }),
                ],
              }),
            ],
          }),
        ],
      }),
    ],
  });
}

function authBadge(required, role = null) {
  const text = required
    ? role
      ? `🔒 Auth required — Role: ${role}`
      : "🔒 Auth required"
    : "🌐 Public";
  const color = required ? "FFF3CD" : "D4EDDA";
  const textColor = required ? "856404" : "155724";
  return new Table({
    width: { size: 3000, type: WidthType.DXA },
    columnWidths: [3000],
    borders: noBorders,
    rows: [
      new TableRow({
        children: [
          new TableCell({
            borders,
            width: { size: 3000, type: WidthType.DXA },
            shading: { fill: color, type: ShadingType.CLEAR },
            margins: { top: 40, bottom: 40, left: 120, right: 120 },
            children: [
              new Paragraph({
                children: [
                  new TextRun({
                    text,
                    font: "Arial",
                    size: 18,
                    color: textColor,
                  }),
                ],
              }),
            ],
          }),
        ],
      }),
    ],
  });
}

function jsonBlock(obj) {
  const lines = JSON.stringify(obj, null, 2).split("\n");
  return [
    new Table({
      width: { size: 9360, type: WidthType.DXA },
      columnWidths: [9360],
      rows: [
        new TableRow({
          children: [
            new TableCell({
              borders,
              width: { size: 9360, type: WidthType.DXA },
              shading: { fill: COLORS.codeBg, type: ShadingType.CLEAR },
              margins: { top: 80, bottom: 80, left: 200, right: 200 },
              children: lines.map(
                (line) =>
                  new Paragraph({
                    spacing: { before: 0, after: 0 },
                    children: [
                      new TextRun({
                        text: line,
                        font: "Courier New",
                        size: 18,
                        color: COLORS.text,
                      }),
                    ],
                  }),
              ),
            }),
          ],
        }),
      ],
    }),
  ];
}

function fieldTable(fields) {
  const headerColor = COLORS.headerBg;
  return new Table({
    width: { size: 9360, type: WidthType.DXA },
    columnWidths: [2200, 1600, 1400, 4160],
    rows: [
      new TableRow({
        children: ["Field", "Type", "Required", "Description"].map((h, i) => {
          const widths = [2200, 1600, 1400, 4160];
          return new TableCell({
            borders,
            width: { size: widths[i], type: WidthType.DXA },
            shading: { fill: headerColor, type: ShadingType.CLEAR },
            margins: { top: 60, bottom: 60, left: 120, right: 120 },
            children: [
              new Paragraph({
                children: [
                  new TextRun({
                    text: h,
                    font: "Arial",
                    size: 18,
                    bold: true,
                    color: "FFFFFF",
                  }),
                ],
              }),
            ],
          });
        }),
      }),
      ...fields.map(
        (f, idx) =>
          new TableRow({
            children: [f.name, f.type, f.required ? "Yes" : "No", f.desc].map(
              (val, i) => {
                const widths = [2200, 1600, 1400, 4160];
                const bg = idx % 2 === 0 ? "FFFFFF" : "F8FAFC";
                const textColor =
                  i === 2
                    ? f.required
                      ? COLORS.green
                      : COLORS.muted
                    : COLORS.text;
                return new TableCell({
                  borders,
                  width: { size: widths[i], type: WidthType.DXA },
                  shading: { fill: bg, type: ShadingType.CLEAR },
                  margins: { top: 60, bottom: 60, left: 120, right: 120 },
                  children: [
                    new Paragraph({
                      children: [
                        new TextRun({
                          text: val,
                          font: i === 0 ? "Courier New" : "Arial",
                          size: 18,
                          color: textColor,
                        }),
                      ],
                    }),
                  ],
                });
              },
            ),
          }),
      ),
    ],
  });
}

function divider() {
  return new Paragraph({
    spacing: { before: 160, after: 160 },
    border: {
      bottom: {
        style: BorderStyle.SINGLE,
        size: 1,
        color: COLORS.border,
        space: 1,
      },
    },
    children: [new TextRun("")],
  });
}

// ── Route Builder ─────────────────────────────────────────────
function route({
  method,
  path,
  description,
  auth,
  role,
  requestFields,
  requestExample,
  responseExample,
  notes,
}) {
  const items = [
    spacer(),
    methodBadge(method, path, description),
    new Paragraph({
      spacing: { before: 80, after: 40 },
      children: [new TextRun("")],
    }),
    authBadge(auth, role),
    spacer(),
  ];

  if (notes) {
    items.push(p(`ℹ️  ${notes}`, { italics: true, color: COLORS.muted }));
    items.push(spacer());
  }

  if (requestFields) {
    items.push(h3("Request Body"));
    items.push(fieldTable(requestFields));
    items.push(spacer());
  }

  if (requestExample) {
    items.push(p("Request Example:", { bold: true }));
    items.push(...jsonBlock(requestExample));
    items.push(spacer());
  }

  if (responseExample) {
    items.push(p("Response Example:", { bold: true }));
    items.push(...jsonBlock(responseExample));
  }

  items.push(divider());
  return items;
}

// ── Document ──────────────────────────────────────────────────
const children = [
  // Cover
  new Paragraph({
    spacing: { before: 480, after: 160 },
    children: [
      new TextRun({
        text: "gotemplate API",
        font: "Arial",
        size: 56,
        bold: true,
        color: COLORS.primary,
      }),
    ],
  }),
  new Paragraph({
    spacing: { before: 0, after: 80 },
    children: [
      new TextRun({
        text: "REST API Reference",
        font: "Arial",
        size: 32,
        color: COLORS.muted,
      }),
    ],
  }),
  new Paragraph({
    spacing: { before: 0, after: 480 },
    children: [
      new TextRun({
        text: `Base URL: http://localhost:8080`,
        font: "Courier New",
        size: 22,
        color: COLORS.accent,
      }),
    ],
  }),
  divider(),

  // Overview
  h1("Overview"),
  p(
    "All API routes are prefixed with /api. Webhooks are at /webhooks. Authentication uses JWT Bearer tokens.",
  ),
  spacer(),
  ...jsonBlock({ Authorization: "Bearer <access_token>" }),
  spacer(),
  p("All responses follow a consistent envelope:"),
  ...jsonBlock({ data: "...", error: "string (on failure)" }),
  divider(),

  // ── AUTH ──────────────────────────────────────────────────
  h1("Authentication"),
  muted("Signup, login, token refresh, and logout."),

  ...route({
    method: "POST",
    path: "/api/auth/signup",
    description: "Register a new user account",
    auth: false,
    requestFields: [
      { name: "name", type: "string", required: true, desc: "Full name" },
      {
        name: "email",
        type: "string",
        required: true,
        desc: "Valid email address",
      },
      {
        name: "password",
        type: "string",
        required: true,
        desc: "Minimum 6 characters",
      },
    ],
    requestExample: {
      name: "John Doe",
      email: "john@example.com",
      password: "secret123",
    },
    responseExample: {
      data: { id: 1, name: "John Doe", email: "john@example.com" },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/auth/login",
    description: "Login and receive JWT tokens",
    auth: false,
    requestFields: [
      {
        name: "email",
        type: "string",
        required: true,
        desc: "Registered email",
      },
      {
        name: "password",
        type: "string",
        required: true,
        desc: "Account password",
      },
    ],
    requestExample: { email: "john@example.com", password: "secret123" },
    responseExample: {
      data: {
        access_token: "eyJ...",
        refresh_token: "eyJ...",
        expires_at: "2026-02-22T15:04:05Z",
      },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/auth/refresh",
    description: "Get new access token using refresh token",
    auth: false,
    requestFields: [
      {
        name: "refresh_token",
        type: "string",
        required: true,
        desc: "Valid refresh token",
      },
    ],
    requestExample: { refresh_token: "eyJ..." },
    responseExample: {
      data: {
        access_token: "eyJ...",
        refresh_token: "eyJ...",
        expires_at: "2026-02-22T15:04:05Z",
      },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/auth/logout",
    description: "Revoke tokens and end session",
    auth: true,
    requestFields: [
      {
        name: "refresh_token",
        type: "string",
        required: true,
        desc: "Refresh token to revoke",
      },
    ],
    requestExample: { refresh_token: "eyJ..." },
    responseExample: { message: "logged out successfully" },
  }),

  // ── USERS ─────────────────────────────────────────────────
  h1("Users"),
  muted("User profile management. Admin-only for list/create/update/delete."),

  ...route({
    method: "GET",
    path: "/api/users/me",
    description: "Get currently authenticated user",
    auth: true,
    responseExample: {
      data: {
        id: 1,
        name: "John Doe",
        email: "john@example.com",
        role: "user",
        is_active: true,
        created_at: "2026-01-01 12:00:00",
      },
    },
  }),

  ...route({
    method: "GET",
    path: "/api/users",
    description: "List all users (paginated)",
    auth: true,
    role: "admin",
    notes: "Query params: ?page=1&limit=10",
    responseExample: {
      data: [
        {
          id: 1,
          name: "John Doe",
          email: "john@example.com",
          role: "user",
          is_active: true,
        },
      ],
      total: 50,
      page: 1,
      limit: 10,
    },
  }),

  ...route({
    method: "GET",
    path: "/api/users/:id",
    description: "Get user by ID",
    auth: true,
    role: "admin",
    responseExample: {
      data: {
        id: 1,
        name: "John Doe",
        email: "john@example.com",
        role: "user",
        is_active: true,
      },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/users",
    description: "Create a new user",
    auth: true,
    role: "admin",
    requestFields: [
      { name: "name", type: "string", required: true, desc: "Full name" },
      { name: "email", type: "string", required: true, desc: "Valid email" },
      {
        name: "password",
        type: "string",
        required: true,
        desc: "Initial password",
      },
    ],
    requestExample: {
      name: "Jane Doe",
      email: "jane@example.com",
      password: "secret123",
    },
    responseExample: {
      data: {
        id: 2,
        name: "Jane Doe",
        email: "jane@example.com",
        role: "user",
        is_active: true,
      },
    },
  }),

  ...route({
    method: "PUT",
    path: "/api/users/:id",
    description: "Update user name or email",
    auth: true,
    role: "admin",
    requestFields: [
      { name: "name", type: "string", required: false, desc: "New name" },
      { name: "email", type: "string", required: false, desc: "New email" },
    ],
    requestExample: { name: "John Updated", email: "updated@example.com" },
    responseExample: {
      data: { id: 1, name: "John Updated", email: "updated@example.com" },
    },
  }),

  ...route({
    method: "DELETE",
    path: "/api/users/:id",
    description: "Delete a user",
    auth: true,
    role: "admin",
    responseExample: { message: "user deleted successfully" },
  }),

  // ── BLOGS ─────────────────────────────────────────────────
  h1("Blogs"),
  muted("Public read access. Write operations require authentication."),

  ...route({
    method: "GET",
    path: "/api/blogs",
    description: "List all published blogs (paginated)",
    auth: false,
    notes: "Query params: ?page=1&limit=10&search=keyword",
    responseExample: {
      data: [
        {
          id: 1,
          title: "Hello World",
          slug: "hello-world",
          status: "published",
          author_id: 1,
          tags: "go,api",
          created_at: "2026-01-01 12:00:00",
        },
      ],
      total: 20,
      page: 1,
      limit: 10,
    },
  }),

  ...route({
    method: "GET",
    path: "/api/blogs/:id",
    description: "Get blog by ID",
    auth: false,
    responseExample: {
      data: {
        id: 1,
        title: "Hello World",
        slug: "hello-world",
        content: "...",
        status: "published",
        author_id: 1,
      },
    },
  }),

  ...route({
    method: "GET",
    path: "/api/blogs/slug/:slug",
    description: "Get blog by slug",
    auth: false,
    responseExample: {
      data: {
        id: 1,
        title: "Hello World",
        slug: "hello-world",
        content: "...",
        status: "published",
      },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/blogs",
    description: "Create a new blog post",
    auth: true,
    requestFields: [
      {
        name: "title",
        type: "string",
        required: true,
        desc: "Blog title (slug auto-generated)",
      },
      {
        name: "content",
        type: "string",
        required: true,
        desc: "Blog body content",
      },
      {
        name: "tags",
        type: "string",
        required: false,
        desc: "Comma-separated tags",
      },
    ],
    requestExample: {
      title: "My First Post",
      content: "Hello world content...",
      tags: "go,backend",
    },
    responseExample: {
      data: {
        id: 1,
        title: "My First Post",
        slug: "my-first-post",
        status: "draft",
        author_id: 1,
      },
    },
  }),

  ...route({
    method: "PUT",
    path: "/api/blogs/:id",
    description: "Update blog post",
    auth: true,
    requestFields: [
      {
        name: "title",
        type: "string",
        required: false,
        desc: "New title (re-generates slug)",
      },
      { name: "content", type: "string", required: false, desc: "New content" },
      {
        name: "status",
        type: "string",
        required: false,
        desc: "draft | published | archived",
      },
      {
        name: "tags",
        type: "string",
        required: false,
        desc: "Comma-separated tags",
      },
    ],
    requestExample: { status: "published", title: "Updated Title" },
    responseExample: {
      data: {
        id: 1,
        title: "Updated Title",
        slug: "updated-title",
        status: "published",
      },
    },
  }),

  ...route({
    method: "DELETE",
    path: "/api/blogs/:id",
    description: "Delete a blog post",
    auth: true,
    responseExample: { message: "blog deleted successfully" },
  }),

  // ── PRODUCTS ──────────────────────────────────────────────
  h1("Products"),
  muted(
    "Public read. Admin-only write. Supports physical and digital product types.",
  ),

  ...route({
    method: "GET",
    path: "/api/products",
    description: "List all active products (paginated)",
    auth: false,
    notes: "Query params: ?page=1&limit=10",
    responseExample: {
      data: [
        {
          id: 1,
          name: "Go Course",
          type: "digital",
          price: 4999,
          currency: "USD",
          is_active: true,
        },
      ],
      total: 5,
      page: 1,
      limit: 10,
    },
  }),

  ...route({
    method: "GET",
    path: "/api/products/:id",
    description: "Get product by ID",
    auth: false,
    responseExample: {
      data: {
        id: 1,
        name: "Go Course",
        description: "Learn Go from scratch",
        type: "digital",
        price: 4999,
        currency: "USD",
        is_active: true,
      },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/products",
    description: "Create a new product",
    auth: true,
    role: "admin",
    requestFields: [
      { name: "name", type: "string", required: true, desc: "Product name" },
      {
        name: "description",
        type: "string",
        required: false,
        desc: "Product description",
      },
      {
        name: "type",
        type: "string",
        required: true,
        desc: "physical | digital",
      },
      {
        name: "price",
        type: "int64",
        required: true,
        desc: "Price in smallest currency unit (cents)",
      },
      {
        name: "currency",
        type: "string",
        required: true,
        desc: "USD | EUR | GBP | INR | JPY | AUD | CAD | SGD | AED",
      },
    ],
    requestExample: {
      name: "Go Course",
      description: "Learn Go from scratch",
      type: "digital",
      price: 4999,
      currency: "USD",
    },
    responseExample: {
      data: {
        id: 1,
        name: "Go Course",
        type: "digital",
        price: 4999,
        currency: "USD",
        is_active: true,
      },
    },
  }),

  ...route({
    method: "PUT",
    path: "/api/products/:id",
    description: "Update a product",
    auth: true,
    role: "admin",
    requestFields: [
      { name: "name", type: "string", required: false, desc: "New name" },
      {
        name: "price",
        type: "int64",
        required: false,
        desc: "New price in cents",
      },
      {
        name: "is_active",
        type: "bool",
        required: false,
        desc: "Enable or disable product",
      },
    ],
    requestExample: { price: 3999, is_active: true },
    responseExample: {
      data: {
        id: 1,
        name: "Go Course",
        price: 3999,
        currency: "USD",
        is_active: true,
      },
    },
  }),

  ...route({
    method: "DELETE",
    path: "/api/products/:id",
    description: "Delete a product",
    auth: true,
    role: "admin",
    responseExample: { message: "product deleted" },
  }),

  // ── ORDERS ────────────────────────────────────────────────
  h1("Orders"),
  muted(
    "Create and track orders. Each order can contain multiple products (batch dispatch).",
  ),

  ...route({
    method: "POST",
    path: "/api/orders",
    description: "Create a new order with one or more products",
    auth: true,
    requestFields: [
      {
        name: "currency",
        type: "string",
        required: true,
        desc: "Order currency (USD, EUR, etc.)",
      },
      {
        name: "items",
        type: "array",
        required: true,
        desc: "Array of order items",
      },
      {
        name: "items[].product_id",
        type: "uint",
        required: true,
        desc: "Product ID",
      },
      {
        name: "items[].quantity",
        type: "int",
        required: true,
        desc: "Quantity (must be > 0)",
      },
    ],
    requestExample: {
      currency: "USD",
      items: [
        { product_id: 1, quantity: 2 },
        { product_id: 3, quantity: 1 },
      ],
    },
    responseExample: {
      data: {
        id: 1,
        user_id: 1,
        status: "pending",
        total_amount: 14997,
        currency: "USD",
        items: [
          { product_id: 1, quantity: 2, unit_price: 4999 },
          { product_id: 3, quantity: 1, unit_price: 4999 },
        ],
      },
    },
  }),

  ...route({
    method: "GET",
    path: "/api/orders/me",
    description: "Get current user's orders",
    auth: true,
    notes: "Query params: ?page=1&limit=10",
    responseExample: {
      data: [
        { id: 1, status: "confirmed", total_amount: 14997, currency: "USD" },
      ],
      total: 3,
      page: 1,
      limit: 10,
    },
  }),

  ...route({
    method: "GET",
    path: "/api/orders/:id",
    description: "Get order by ID",
    auth: true,
    responseExample: {
      data: {
        id: 1,
        user_id: 1,
        status: "dispatched",
        total_amount: 14997,
        items: [],
      },
    },
  }),

  ...route({
    method: "PATCH",
    path: "/api/orders/:id/status",
    description: "Update order status (admin only)",
    auth: true,
    role: "admin",
    notes:
      "Valid transitions: pending→confirmed, confirmed→dispatched, dispatched→delivered, delivered→completed. pending/confirmed→cancelled.",
    requestFields: [
      {
        name: "status",
        type: "string",
        required: true,
        desc: "pending | confirmed | dispatched | delivered | completed | cancelled",
      },
    ],
    requestExample: { status: "dispatched" },
    responseExample: { message: "order status updated" },
  }),

  // ── PAYMENTS ──────────────────────────────────────────────
  h1("Payments"),
  muted(
    "Payment history for compliance. Records every transaction with full audit trail.",
  ),

  ...route({
    method: "GET",
    path: "/api/payments/:id",
    description: "Get payment record by ID",
    auth: true,
    responseExample: {
      data: {
        id: "pay_1234567890",
        user_id: 1,
        order_id: "ord_123",
        amount: 4999,
        currency: "USD",
        status: "completed",
        provider: "stripe",
        provider_txn_id: "pi_xxx",
        customer_email: "john@example.com",
        created_at: "2026-01-01T12:00:00Z",
      },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/payments/:id/refund",
    description: "Refund a completed payment",
    auth: true,
    role: "admin",
    notes: 'Only payments with status "completed" can be refunded.',
    responseExample: { message: "payment refunded" },
  }),

  ...route({
    method: "POST",
    path: "/api/payments/:id/complete",
    description: "Mark payment as completed",
    auth: true,
    role: "admin",
    notes: "Typically called internally via webhook processing.",
    responseExample: { message: "payment marked completed" },
  }),

  ...route({
    method: "POST",
    path: "/api/payments/:id/fail",
    description: "Mark payment as failed",
    auth: true,
    role: "admin",
    responseExample: { message: "payment marked failed" },
  }),

  // ── MEMBERSHIPS ───────────────────────────────────────────
  h1("Memberships"),
  muted("User subscription plans. Plans: free, basic, pro, enterprise."),

  ...route({
    method: "GET",
    path: "/api/memberships/me",
    description: "Get current user's membership",
    auth: true,
    responseExample: {
      data: {
        id: 1,
        user_id: 1,
        plan: "pro",
        status: "active",
        starts_at: "2026-01-01T00:00:00Z",
        expires_at: "2026-02-01T00:00:00Z",
      },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/memberships/activate",
    description: "Activate a membership plan",
    auth: true,
    notes: "User must not already have an active membership.",
    requestFields: [
      {
        name: "plan",
        type: "string",
        required: true,
        desc: "free | basic | pro | enterprise",
      },
    ],
    requestExample: { plan: "pro" },
    responseExample: {
      data: {
        id: 1,
        plan: "pro",
        status: "active",
        expires_at: "2026-02-22T00:00:00Z",
      },
    },
  }),

  ...route({
    method: "POST",
    path: "/api/memberships/upgrade",
    description: "Upgrade to a higher plan",
    auth: true,
    notes: "Can only upgrade to a higher tier. free→basic→pro→enterprise.",
    requestFields: [
      {
        name: "plan",
        type: "string",
        required: true,
        desc: "Target plan (must be higher than current)",
      },
    ],
    requestExample: { plan: "enterprise" },
    responseExample: {
      data: {
        id: 1,
        plan: "enterprise",
        status: "active",
        expires_at: "2027-02-22T00:00:00Z",
      },
    },
  }),

  ...route({
    method: "DELETE",
    path: "/api/memberships/cancel",
    description: "Cancel active membership",
    auth: true,
    responseExample: { message: "membership cancelled" },
  }),

  // ── WEBHOOKS ──────────────────────────────────────────────
  h1("Webhooks"),
  muted(
    "Public endpoints. Signature is verified server-side using HMAC-SHA256.",
  ),

  ...route({
    method: "POST",
    path: "/webhooks/stripe",
    description: "Receive Stripe webhook events",
    auth: false,
    notes:
      "Signature verified via X-Stripe-Signature header using HMAC-SHA256. Handles: payment_intent.succeeded, payment_intent.payment_failed, charge.refunded.",
    requestFields: [
      { name: "id", type: "string", required: true, desc: "Stripe event ID" },
      {
        name: "type",
        type: "string",
        required: true,
        desc: "Event type (e.g. payment_intent.succeeded)",
      },
      { name: "data", type: "object", required: true, desc: "Event payload" },
    ],
    responseExample: { received: true },
  }),

  ...route({
    method: "POST",
    path: "/webhooks/razorpay",
    description: "Receive Razorpay webhook events",
    auth: false,
    notes:
      "Signature verified via X-Razorpay-Signature header using HMAC-SHA256. Handles: payment.captured, payment.failed, refund.processed.",
    requestFields: [
      { name: "id", type: "string", required: true, desc: "Razorpay event ID" },
      {
        name: "event",
        type: "string",
        required: true,
        desc: "Event type (e.g. payment.captured)",
      },
      {
        name: "payload",
        type: "object",
        required: true,
        desc: "Event payload",
      },
    ],
    responseExample: { received: true },
  }),

  // ── ERROR REFERENCE ───────────────────────────────────────
  h1("Error Reference"),
  spacer(),
  fieldTable([
    {
      name: "400",
      type: "Bad Request",
      required: false,
      desc: "Invalid request body or missing required fields",
    },
    {
      name: "401",
      type: "Unauthorized",
      required: false,
      desc: "Missing or invalid JWT token",
    },
    {
      name: "403",
      type: "Forbidden",
      required: false,
      desc: "Insufficient role permissions",
    },
    {
      name: "404",
      type: "Not Found",
      required: false,
      desc: "Resource not found",
    },
    {
      name: "500",
      type: "Server Error",
      required: false,
      desc: "Unexpected internal error",
    },
  ]),
  spacer(),
  p("All errors return the following shape:"),
  ...jsonBlock({ error: "description of what went wrong" }),
];

const doc = new Document({
  styles: {
    default: {
      document: { run: { font: "Arial", size: 22, color: COLORS.text } },
    },
    paragraphStyles: [
      {
        id: "Heading1",
        name: "Heading 1",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 36, bold: true, font: "Arial", color: COLORS.primary },
        paragraph: { spacing: { before: 360, after: 120 }, outlineLevel: 0 },
      },
      {
        id: "Heading2",
        name: "Heading 2",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 28, bold: true, font: "Arial", color: COLORS.primary },
        paragraph: { spacing: { before: 280, after: 80 }, outlineLevel: 1 },
      },
      {
        id: "Heading3",
        name: "Heading 3",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 24, bold: true, font: "Arial", color: COLORS.accent },
        paragraph: { spacing: { before: 200, after: 60 }, outlineLevel: 2 },
      },
    ],
  },
  sections: [
    {
      properties: {
        page: {
          size: { width: 12240, height: 15840 },
          margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 },
        },
      },
      footers: {
        default: new Footer({
          children: [
            new Paragraph({
              alignment: AlignmentType.CENTER,
              children: [
                new TextRun({
                  text: "gotemplate API Reference  •  Page ",
                  font: "Arial",
                  size: 18,
                  color: COLORS.muted,
                }),
                new TextRun({
                  children: [PageNumber.CURRENT],
                  font: "Arial",
                  size: 18,
                  color: COLORS.muted,
                }),
              ],
            }),
          ],
        }),
      },
      children,
    },
  ],
});

Packer.toBuffer(doc).then((buffer) => {
  fs.writeFileSync("/mnt/user-data/outputs/api-docs.docx", buffer);
  console.log("Done: api-docs.docx");
});
