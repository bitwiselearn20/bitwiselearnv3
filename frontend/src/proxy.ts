import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { match } from "path-to-regexp";

import { URL_ACCESS_MAP } from "./lib/access";
import { checkJWT } from "./lib/authJwt";

/**
 * Matches a pathname against role-based route patterns
 */
function isRouteAllowed(pathname: string, routes: string[]) {
  return routes.some((route) => {
    const matcher = match(route, {
      decode: decodeURIComponent,
      end: true, // exact match, no extra segments
    });
    return matcher(pathname);
  });
}

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;

  /**
   * Public routes
   */
  const PUBLIC_ROUTES = [
    "/",
    "/about",
    "/contact",
    "/student-login",
    "/multi-login",
    "/admin-login",
    "/api/run",
    "/our-services",
    "/blog",
    "/listed-courses",
    "/services",
    "/blog",
    "/blog/:id",
  ];

  if (PUBLIC_ROUTES.includes(pathname)) {
    return NextResponse.next();
  }

  /**
   * Auth check
   */
  const token = request.cookies.get("token")?.value;
  let role: string | undefined;

  if (token) {
    try {
      const jwtData = checkJWT(token);
      role = String(jwtData?.type || "").toUpperCase();
    } catch {
      role = undefined;
    }
  }

  if (!role) {
    return NextResponse.redirect(new URL("/student-login", request.url));
  }
  if (role && role.includes("LOGIN")) {
    switch (role) {
      case "ADMIN":
      case "SUPERADMIN":
        return NextResponse.redirect(new URL("/admin-dashboard", request.url));

      case "TEACHER":
        return NextResponse.redirect(
          new URL("/teacher-dashboard", request.url),
        );

      case "STUDENT":
        return NextResponse.redirect(new URL("/dashboard", request.url));

      case "INSTITUTION":
        return NextResponse.redirect(
          new URL("/institution-dashboard", request.url),
        );

      case "VENDOR":
        return NextResponse.redirect(new URL("/vendor-dashboard", request.url));

      default:
        return NextResponse.redirect(new URL("/", request.url));
    }
  }

  /**
   * Resolve allowed routes
   */
  let allowedRoutes = URL_ACCESS_MAP[role];

  // SUPERADMIN has its own dedicated routes
  if (role === "SUPERADMIN" && !allowedRoutes) {
    allowedRoutes = URL_ACCESS_MAP["ADMIN"];
  }

  if (!allowedRoutes) {
    return NextResponse.redirect(new URL("/student-login", request.url));
  }

  /**
   * Authorization check
   */
  if (isRouteAllowed(pathname, allowedRoutes)) {
    return NextResponse.next();
  }

  /**
   * Fallback redirects per role
   */
  switch (role) {
    case "ADMIN":
    case "SUPERADMIN":
      return NextResponse.redirect(new URL("/admin-dashboard", request.url));

    case "TEACHER":
      return NextResponse.redirect(new URL("/teacher-dashboard", request.url));

    case "STUDENT":
      return NextResponse.redirect(new URL("/dashboard", request.url));

    case "INSTITUTION":
      return NextResponse.redirect(
        new URL("/institution-dashboard", request.url),
      );

    case "VENDOR":
      return NextResponse.redirect(new URL("/vendor-dashboard", request.url));

    default:
      return NextResponse.redirect(new URL("/", request.url));
  }
}

export const config = {
  matcher: ["/((?!api|_next|favicon.ico).*)"],
};
