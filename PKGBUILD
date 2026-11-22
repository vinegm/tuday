# Maintainer: vinegm <info AT vinegm D0T dev>
pkgname='tuday-bin'
pkgver=1.0.0
pkgrel=1
pkgdesc="Doomsday TUI trainer"
arch=('x86_64' 'aarch64')
url="https://github.com/vinegm/tuday"
license=('MIT')
provides=('tuday')
conflicts=('tuday')
options=('!strip' '!debug')
source_x86_64=("tuday-${pkgver}-x86_64.tar.gz::${url}/releases/download/v${pkgver}/tuday_${pkgver}_linux_amd64.tar.gz")
source_aarch64=("tuday-${pkgver}-aarch64.tar.gz::${url}/releases/download/v${pkgver}/tuday_${pkgver}_linux_arm64.tar.gz")
sha256sums_x86_64=('9b8db227b687068becb196c58e2e321c467572a90cabab4561e9b6591467e1fe')
sha256sums_aarch64=('297b414352a65848f3dd3acd23ae37f6d05f2e6ba55e5f5886755fd7fdc43108')

package() {
    install -Dm755 tuday "$pkgdir/usr/bin/tuday"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
