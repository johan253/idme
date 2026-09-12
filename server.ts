import express, { Request, Response, NextFunction } from "express";
import cookieParser from "cookie-parser";
import { config } from "dotenv";
import jwt from "jsonwebtoken";
import * as bcrypt from "bcryptjs";

config();
const JWT_SECRET = process.env.JWT_SECRET || "";
const PORT = process.env.PORT || 3000;

const app = express();
app.use(express.json());
app.use(cookieParser());

// data

interface JwtPayload {
  iss: string;
  sub: string;
  iat: number;
  exp: number;
}

interface AuthenticatedRequest extends Request {
  user?: JwtPayload;
}

interface User {
  username: string;
  hashedPassword: string;
}
const users: User[] = [];

interface Order {
  id: number;
  username: string;
  items: string[];
}
const orders: Order[] = [];

// helpers

const userExists = (username: string) => {
  return users.some((user) => user.username === username);
}

const authenticateToken = (req: Request, res: Response, next: NextFunction) => {
  const cookies = req.cookies;
  const token = cookies && cookies.access_token;
  if (!token) {
    return res.status(401).json({ message: "Unauthorized" });
  }
  try {
    const decoded = jwt.verify(token, JWT_SECRET) as JwtPayload;
    (req as AuthenticatedRequest).user = decoded;
    next();
  } catch (err) {
    return res.status(403).json({ message: "Forbidden" });
  }
}

// routes

app.post("/register", (req, res) => {
  const { username, password } = req.body;
  if (!username || !password || userExists(username)) {
    return res.status(400).json({ message: "Invalid username or password" });
  }
  users.push({
    username,
    hashedPassword: bcrypt.hashSync(password, 10),
  })
  return res.status(201).json({ message: "User registered successfully" });
});

app.post("/login", (req, res) => {
  const { username, password } = req.body;
  const user = users.find((user) => user.username === username);
  if (!username || !password || !user || !bcrypt.compareSync(password, user.hashedPassword)) {
    return res.status(400).json({ message: "Invalid username or password" });
  }
  const token = jwt.sign(
    { sub: username },
    JWT_SECRET,
    { expiresIn: "1m", issuer: "login" }
  );
  return res
    .status(200)
    .cookie("access_token", token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
    })
    .json({ message: "Login successful" })
});

app.post("/logout", (_, res) => {
  return res
    .clearCookie("access_token")
    .clearCookie("refresh_token")
    .json({ message: "Logout successful" });
});

app.post("/orders", authenticateToken, (req: AuthenticatedRequest, res: Response) => {
  const username = req.user?.sub;
  if (!username) {
    return res.status(401).json({ message: "Unauthorized" });
  }
  const { items } = req.body;
  if (!items || !Array.isArray(items)) {
    return res.status(400).json({ message: "Invalid items" });
  }
  const order: Order = {
    id: orders.length + 1,
    username,
    items,
  };
  orders.push(order);
  return res.status(201).json({ message: "Order created successfully", order });
});

app.get("/orders", authenticateToken, (req: AuthenticatedRequest, res: Response) => {
  const username = req.user?.sub;
  if (!username) {
    return res.status(401).json({ message: "Unauthorized" });
  }
  const userOrders = orders.filter((order) => order.username === username);
  return res.status(200).json({ orders: userOrders });
});

app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
