#include "src/stirling/obj_tools/obj_tools.h"

#include "src/common/testing/test_environment.h"
#include "src/common/testing/testing.h"

DEFINE_string(go_grpc_client_path, "", "The path to the go greeter client executable.");

namespace pl {
namespace stirling {
namespace obj_tools {

using ::testing::_;
using ::testing::Contains;
using ::testing::ElementsAre;
using ::testing::EndsWith;
using ::testing::IsEmpty;
using ::testing::Pair;

// Tests GetActiveBinaries() resolves this running test itself.
// We instruct GetActiveBinaries() to behave as if the test is not running inside a container.
TEST(GetActiveBinariesTest, CaptureTestBinary) {
  int32_t mypid = getpid();
  std::map<int32_t, std::filesystem::path> pid_paths = {
      {mypid, std::filesystem::path("/proc") / std::to_string(mypid)}};
  const std::map<std::string, std::vector<int>> binaries = GetActiveBinaries(pid_paths);
  EXPECT_THAT(binaries, Contains(Pair(EndsWith("src/stirling/obj_tools/obj_tools_test"), _)))
      << "Should see the test process itself";
}

TEST(GetActiveBinaryTest, CaptureTestBinary) {
  auto binary_or =
      GetActiveBinary(/*host_path*/ {}, std::filesystem::path("/proc") / std::to_string(getpid()));
  EXPECT_OK_AND_THAT(binary_or, EndsWith("src/stirling/obj_tools/obj_tools_test"));
}

// TEST(GetSymAddrsTest, SymbolAddress) {
//  CHECK(!FLAGS_go_grpc_client_path.empty())
//      << "--go_grpc_client_path cannot be empty. You should run this test with bazel.";
//  std::map<std::string, std::vector<int>> binaries;
//  binaries[FLAGS_go_grpc_client_path] = {1, 2, 3};
//  EXPECT_THAT(GetSymAddrs(binaries), ElementsAre(Pair(1, _), Pair(2, _), Pair(3, _)));
//}

TEST(ObjToolsTest, FindRetInsts) {
  // Note \xc2 \xca must be followed by 2 bytes, it cannot be arbitrary values. For example,
  // \xc2\x00\x00 is rejected by disassembler.
  constexpr char kData[] = "\x55\xc3\xcb\xc2\x01\x01\xca\x01\x01\xc3\xcb";

  utils::u8string_view byte_codes(reinterpret_cast<const uint8_t*>(kData));

  LLVMDisasmContext disasm_ctx;

  EXPECT_THAT(FindRetInsts(disasm_ctx, byte_codes), ElementsAre(1, 2, 3, 6, 9, 10));
}

}  // namespace obj_tools
}  // namespace stirling
}  // namespace pl
