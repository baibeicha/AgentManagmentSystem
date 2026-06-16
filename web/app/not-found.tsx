import Link from 'next/link';

export default function NotFound() {
  return (
    <div className="flex h-full flex-col items-center justify-center space-y-4">
      <h2 className="text-xl font-bold">Page Not Found</h2>
      <Link href="/" className="text-sky-400 hover:text-sky-300">
        Return Home
      </Link>
    </div>
  );
}
