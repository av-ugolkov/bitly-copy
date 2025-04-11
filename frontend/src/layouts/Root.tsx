import { Outlet, useNavigate } from 'react-router';

export default function BaseLayout() {
  const navigate = useNavigate();

  return (
    <>
      <header className='flex flex-col min-w-[640px] pt-6 gap-y-2'>
        <h1 className='flex justify-center text-4xl'>Bitly Copy</h1>
        <div className='flex justify-center'>
          {Button('Gamerate', () => navigate('/'))}
          {Button('Get URLs', () => navigate('/get-urls'))}
          {Button('Statistics', () => navigate('/statistics'))}
        </div>
      </header>
      <main className='flex justify-center'>
        <Outlet />
      </main>
    </>
  );
}

function Button(text: string, onClick: () => void) {
  return (
    <button
      className='bg-gray-200 hover:bg-gray-300 p-2 m-2 rounded-lg'
      onClick={onClick}>
      {text}
    </button>
  );
}
