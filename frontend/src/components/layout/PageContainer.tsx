import { ReactNode } from "react";

interface Props {
  title: string;
  description?: string;
  children: ReactNode;
  action?: ReactNode;
}

export function PageContainer({
  title,
  description,
  children,
  action
}: Props) {
  return (
    <div>
      <div className="mb-8 flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold">
            {title}
          </h1>

          {description && (
            <p className="mt-2 text-zinc-400">
              {description}
            </p>
          )}
        </div>

        {action}
      </div>

      {children}
    </div>
  );
}