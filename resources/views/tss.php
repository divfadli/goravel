<!DOCTYPE html>
<html lang="id">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- <meta name="csrf-token" content="Ay2WXRnDbngaffmDTZQcWJRXYs8WRi1rKpCQKFgq"> -->

    <title>Rekapitulasi - Login </title>

    <!-- Scripts -->
    <link rel="preload" as="style" href="https://laratest.nataon.dev/build/assets/app-012cb4dc.css" />
    <link rel="modulepreload" href="https://laratest.nataon.dev/build/assets/app-6f64ce1f.js" />
    <link rel="stylesheet" href="https://laratest.nataon.dev/build/assets/app-012cb4dc.css" data-navigate-track="reload" />
    <script type="module" src="https://laratest.nataon.dev/build/assets/app-6f64ce1f.js" data-navigate-track="reload"></script><!-- Livewire Styles -->
    <style>
        [wire\:loading][wire\:loading],
        [wire\:loading\.delay][wire\:loading\.delay],
        [wire\:loading\.inline-block][wire\:loading\.inline-block],
        [wire\:loading\.inline][wire\:loading\.inline],
        [wire\:loading\.block][wire\:loading\.block],
        [wire\:loading\.flex][wire\:loading\.flex],
        [wire\:loading\.table][wire\:loading\.table],
        [wire\:loading\.grid][wire\:loading\.grid],
        [wire\:loading\.inline-flex][wire\:loading\.inline-flex] {
            display: none;
        }

        [wire\:loading\.delay\.none][wire\:loading\.delay\.none],
        [wire\:loading\.delay\.shortest][wire\:loading\.delay\.shortest],
        [wire\:loading\.delay\.shorter][wire\:loading\.delay\.shorter],
        [wire\:loading\.delay\.short][wire\:loading\.delay\.short],
        [wire\:loading\.delay\.default][wire\:loading\.delay\.default],
        [wire\:loading\.delay\.long][wire\:loading\.delay\.long],
        [wire\:loading\.delay\.longer][wire\:loading\.delay\.longer],
        [wire\:loading\.delay\.longest][wire\:loading\.delay\.longest] {
            display: none;
        }

        [wire\:offline][wire\:offline] {
            display: none;
        }

        [wire\:dirty]:not(textarea):not(input):not(select) {
            display: none;
        }

        :root {
            --livewire-progress-bar-color: #2299dd;
        }

        [x-cloak] {
            display: none !important;
        }
    </style>
</head>

<body class="font-sans antialiased bg-gray-50/90 dark:bg-bluegray">

    <div x-persist="background">
        <ul class="background-animate-cube">
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
        </ul>
    </div>

    <div class="flex flex-col items-center min-h-screen py-16 sm:justify-center sm:pt-5">
        <div class="w-full max-w-md p-6 mx-auto mt-6 overflow-hidden card">
            <div class="text-center">
                <div class="bg-gradient-to-br bg-clip-text box-decoration-clone from-blue-400 to-blue-700 font-bold text-transparent !no-underline text-xl lowercase font-poppins">
                    <a href="/">
                        Rekapitulasi
                    </a>
                </div>
                <!-- <div class="mt-2 mb-3 text-sm text-gray-600 dark:text-gray-400">
                    Belum punya akun?
                    <a href="https://laratest.nataon.dev/register" class="text-blue-600 hover:underline" wire:navigate.hover>
                        Daftar
                    </a>
                </div> -->
            </div>
            <div wire:snapshot="{&quot;data&quot;:{&quot;form&quot;:[{&quot;username&quot;:&quot;&quot;,&quot;password&quot;:&quot;&quot;,&quot;remember&quot;:false},{&quot;class&quot;:&quot;App\\Livewire\\Forms\\LoginForm&quot;,&quot;s&quot;:&quot;form&quot;}]},&quot;memo&quot;:{&quot;id&quot;:&quot;5ZLNGbFRumWfFS1bbc7S&quot;,&quot;name&quot;:&quot;pages.auth.login&quot;,&quot;path&quot;:&quot;login&quot;,&quot;method&quot;:&quot;GET&quot;,&quot;children&quot;:[],&quot;scripts&quot;:[],&quot;assets&quot;:[],&quot;errors&quot;:[],&quot;locale&quot;:&quot;id&quot;},&quot;checksum&quot;:&quot;9e826023b9d5a5e63cbebb1af875b3daf09b09cca843281080866153409e624f&quot;}" wire:effects="[]" wire:id="5ZLNGbFRumWfFS1bbc7S">
                <form wire:submit="login">
                    <div class="grid gap-y-4">

                        <!-- Username or Email -->
                        <div>
                            <label class="block text-sm mb-2 dark:text-white" for="username">
                                Username
                            </label>
                            <input class="py-3 px-4 block w-full border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-bluegray dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600 block w-full mt-1" wire:model="form.username" id="username" type="text" name="username" required="required" autofocus="autofocus" autocomplete="username" placeholder="Username">
                            <!--[if BLOCK]><![endif]--> <!--[if ENDBLOCK]><![endif]-->
                        </div>

                        <!-- Password -->
                        <div>
                            <div class="flex items-center justify-between">
                                <label class="block text-sm mb-2 dark:text-white" for="password">
                                    Kata Sandi
                                </label>
                            </div>

                            <input class="py-3 px-4 block w-full border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-bluegray dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600 block w-full" wire:model="form.password" id="password" type="password" name="password" required="required" autocomplete="current-password" placeholder="Password">

                            <!--[if BLOCK]><![endif]--> <!--[if ENDBLOCK]><![endif]-->
                        </div>

                        <!-- Remember Me -->
                        <div class="block">
                            <label for="remember" class="inline-flex items-center">
                                <input wire:model="form.remember" id="remember" type="checkbox" class="text-blue-600 border-gray-300 rounded shadow-sm dark:bg-gray-900 dark:border-gray-700 focus:ring-blue-500 dark:focus:ring-blue-600 dark:focus:ring-offset-gray-800" name="remember">
                                <span class="text-sm text-gray-600 ms-2 dark:text-gray-400">Ingat saya</span>
                            </label>
                        </div>
                        <div class="grid">
                            <button type="submit" class="btn info">
                                <!--[if BLOCK]><![endif]--> <span class="animate-spin inline-block w-4 h-4 border-[3px] border-current border-t-transparent text-white rounded-full" role="status" aria-label="loading" wire:loading.delay.long>
                                    <span class="sr-only">Loading...</span>
                                </span>
                                <!--[if ENDBLOCK]><![endif]-->
                                Sign in
                            </button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    </div>
    <div role="status" id="toaster" x-data="toasterHub([], JSON.parse('{\u0022alignment\u0022:\u0022top\u0022,\u0022duration\u0022:3000}'))" class="fixed z-[60] p-4 w-full flex flex-col pointer-events-none sm:p-6 top-0 items-end rtl:items-start">
        <template x-for="toast in toasts" :key="toast.id">
            <div x-show="toast.isVisible" x-init="$nextTick(() => toast.show($el))" x-transition:enter-start="-translate-y-12 opacity-0" x-transition:enter-end="translate-y-0 opacity-100" x-transition:leave-end="opacity-0 scale-90" class="relative duration-300 transform transition ease-in-out max-w-xs w-full pointer-events-auto" :class="toast.select({ error: 'text-white', info: 'text-white', success: 'text-white', warning: 'text-white' })">
                <i x-text="toast.message" class="inline-block select-none not-italic px-6 py-3.5 rounded-xl shadow-lg text-sm w-full mb-3" :class="toast.select({
                    error: 'bg-red-500',
                    info: 'bg-blue-500',
                    success: 'bg-teal-500',
                    warning: 'bg-yellow-500'
                })">
                </i>

                <button @click="toast.dispose()" aria-label="close" class="absolute right-0 px-2 focus:outline-none rtl:right-auto rtl:left-0 top-0">
                    <iconify-icon icon="iconamoon:close-thin" width="20" class="mt-2"></iconify-icon>
                </button>
            </div>
        </template>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/iconify-icon@2.0.0/dist/iconify-icon.min.js"></script>
    <!-- Livewire Scripts -->
    <script src="/livewire/livewire.js?id=f121a5df" data-csrf="Ay2WXRnDbngaffmDTZQcWJRXYs8WRi1rKpCQKFgq" data-update-uri="/livewire/update" data-navigate-once="true"></script>
</body>

</html>