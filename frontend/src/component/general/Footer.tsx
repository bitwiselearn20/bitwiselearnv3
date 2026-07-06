import { Facebook, Instagram, Linkedin, Twitter } from "lucide-react";
import Link from "next/link";
import { getColors } from "@/component/general/(Color Manager)/useColors";
import Image from "next/image";
import logo from "../../../public/images/Logo.png";

const quickLinks = [
  { href: "/", label: "Home" },
  { href: "/about", label: "About" },
  { href: "/listed-courses", label: "Courses" },
  { href: "/our-services", label: "Services" },
  { href: "/contact", label: "Contact Us" },
];

export default function Footer() {
  const Colors = getColors();

  return (
    <footer
      className={`${Colors.background.primary} ${Colors.border.default} backdrop-blur-lg py-10 px-6 transition-colors duration-300`}
    >
      <div className="max-w-7xl mx-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-10">
        {/* Company Info & Social Media */}
        <div className="flex flex-col items-center sm:items-start text-center sm:text-left">
          <Link href="/">
            <Image src={logo} alt="Logo" height={40} />
          </Link>
          <p className={`text-sm max-w-70 mb-4 ${Colors.text.special}`}>
            Bridging the gap between academic theory and industry reality — one placement-ready student at a time.
          </p>
          <div className="flex space-x-4">
            <a
              href="#"
              aria-label="Facebook"
              className={`transition-colors ${Colors.text.primary} hover:text-blue-600`}
            >
              <Facebook size={22} />
            </a>
            <a
              href="#"
              aria-label="Twitter"
              className={`transition-colors ${Colors.text.primary} hover:text-sky-400`}
            >
              <Twitter size={22} />
            </a>
            <a
              href="#"
              aria-label="Instagram"
              className={`transition-colors ${Colors.text.primary} hover:text-yellow-400`}
            >
              <Instagram size={22} />
            </a>
            <a
              href="#"
              aria-label="LinkedIn"
              className={`transition-colors ${Colors.text.primary} hover:text-blue-400`}
            >
              <Linkedin size={22} />
            </a>
          </div>
        </div>

        {/* Quick Links */}
        <div className="text-center sm:text-left">
          <h4 className={`font-bold text-lg mb-4 ${Colors.text.primary}`}>
            Quick Links
          </h4>
          <ul className={`space-y-2 ${Colors.text.secondary}`}>
            {quickLinks.map((link) => (
              <li key={link.href}>
                <Link
                  href={link.href}
                  className="hover:font-semibold transition-all duration-200"
                >
                  {link.label}
                </Link>
              </li>
            ))}
            <li>
              <Link
                href="/multi-login"
                className="hover:font-semibold transition-all duration-200"
              >
                Student / Institute Login
              </Link>
            </li>
          </ul>
        </div>

        {/* Contact Info */}
        <div className="text-center sm:text-left">
          <h4 className={`font-bold text-lg mb-4 ${Colors.text.primary}`}>
            Contact
          </h4>
          <ul className={`space-y-2 ${Colors.text.secondary}`}>
            <li>
              <p>Email: sales_support@bitwiselearn.com</p>
            </li>
            <li>
              <p>Phone: +91 9787777547</p>
            </li>
            <li>Address: Bangalore, India</li>
          </ul>
        </div>
      </div>

      {/* Divider + Copyright + Legal */}
      <div
        className={`mt-8 pt-4 border-t ${Colors.background.primary} flex flex-col items-center justify-between gap-3 text-center text-sm ${Colors.text.secondary} sm:flex-row sm:text-left`}
      >
        <p>&copy; {new Date().getFullYear()} Bitwise Learn. All rights reserved.</p>
        <div className="flex gap-5">
          <a href="#" className="hover:font-semibold transition-all duration-200">Privacy Policy</a>
          <a href="#" className="hover:font-semibold transition-all duration-200">Terms of Service</a>
        </div>
      </div>
    </footer>
  );
}
